package vclient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"v2ray.com/core/transport/internet"
	"v2ray.com/core/watchman/controllers"
	"v2ray.com/core/watchman/database"
	"v2ray.com/core/watchman/logging"
	"v2ray.com/core/watchman/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	DefaultVmessInboundTag = "vmess-proxy"
	DefaultVlessInboundTag = "vless-proxy"
)

type VClient struct {
	Conn            *grpc.ClientConn
	Stats           *StatsServiceClient
	VmessManager    *HandlerServiceClient
	VlessManager    *HandlerServiceClient
	Logger          *zap.Logger
	Accounts        map[string]*proto.UserModel
	VmessInboundTag string
	VlessInboundTag string
}

func Connect(address string, timeoutDuration time.Duration) (*VClient, error) {
	timeout := time.After(timeoutDuration)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return nil, errors.New("connect timeout")
		case <-tick:
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err == nil {
				return &VClient{
					Conn:   conn,
					Logger: logging.GetInstance().Logger,
				}, nil
			}
			return nil, errors.New("connect failed")
		}
	}
}

func (v *VClient) InitServices(vmessInboundTag, vlessInboundTag string) error {
	if vmessInboundTag != "" {
		v.VmessInboundTag = vmessInboundTag
	}

	if vlessInboundTag != "" {
		v.VlessInboundTag = vlessInboundTag
	}

	if v.Conn == nil {
		return errors.New("connect failed")
	}

	if v.Stats == nil {
		v.Logger.Debug("start stats service")
		v.Stats = NewStatsServiceClient(v.Conn)
	}

	if v.VmessManager == nil && v.VmessInboundTag != "" {
		v.Logger.Debug("start vmess manage service")
		v.VmessManager = NewHandlerServiceClient(v.Conn, v.VmessInboundTag, false)
	}

	if v.VlessManager == nil && v.VlessInboundTag != "" {
		v.Logger.Debug("start vless manage service")
		v.VlessManager = NewHandlerServiceClient(v.Conn, v.VlessInboundTag, true)
	}

	return nil
}

func (v *VClient) Startup(dbUrl string, nodeId, checkRate int64) error {
	database.Connect(context.TODO(), v.Logger, &proto.DBConfig{
		Master:  dbUrl,
		MaxIdle: 10,
		MaxOpen: 10,
	})

	if err := v.Sync(nodeId); err != nil {
		v.Logger.Error(err.Error())
	}

	tick := time.Tick(time.Second * time.Duration(checkRate))

	for c := range tick {
		v.Logger.Info("sync database: ", zap.String("time", c.String()))
		if err := v.Sync(nodeId); err != nil {
			v.Logger.Error(err.Error())
		}
	}

	return nil
}

func (v *VClient) AddMainInbound(port uint16) error {
	streamSetting := &internet.StreamConfig{}
	if err := v.VmessManager.AddVmessInbound(port, "0.0.0.0", streamSetting); err != nil {
		return err
	} else {
		v.Logger.Debug(fmt.Sprintf("Successfully add MAIN INBOUND %s port %d", "0.0.0.0", port))
	}
	return nil
}

func (v *VClient) Sync(nodeId int64) error {
	if err := v.syncTraffics(nodeId); err != nil {
		v.Logger.Error(err.Error())
		return err
	}

	v.syncAccounts(nodeId)
	return nil
}

func (v *VClient) loadAccounts(nodeId int64) []*proto.UserModel {
	var accounts []*proto.UserModel
	if err := controllers.Agent.GetAccounts(context.TODO(), nodeId, &accounts); err != nil {
		v.Logger.Error(err.Error())
	}
	return accounts
}

func (v *VClient) syncAccounts(nodeId int64) {
	var addedUsers []*proto.UserModel
	var modifiedUsers []*proto.UserModel
	var removedUsers []*proto.UserModel
	accounts := v.loadAccounts(nodeId)

	newAccounts := make(map[string]*proto.UserModel)

	for i := range accounts {
		user := accounts[i]
		newAccounts[user.Email] = user
		if u, ok := v.Accounts[user.Email]; !ok {
			addedUsers = append(addedUsers, user)
		} else if u.UUID != user.UUID {
			modifiedUsers = append(modifiedUsers, user)
		}
	}

	for k, a := range v.Accounts {
		if _, ok := newAccounts[k]; !ok {
			removedUsers = append(removedUsers, a)
		}
	}

	for i := range addedUsers {
		user := addedUsers[i]
		if v.VmessManager != nil {
			if err := v.VmessManager.AddUser(user); err != nil {
				v.Logger.Error(err.Error())
			}
			v.Logger.Info("新增Vmess用户", zap.String("email", user.Email), zap.String("uuid", user.UUID))
		}

		if v.VlessManager != nil {
			if err := v.VlessManager.AddUser(user); err != nil {
				v.Logger.Error(err.Error())
			}
			v.Logger.Info("新增Vless用户", zap.String("email", user.Email), zap.String("uuid", user.UUID))
		}

	}

	for i := range modifiedUsers {
		user := modifiedUsers[i]
		if v.VmessManager != nil {
			if err := v.VmessManager.DelUser(user.Email); err != nil {
				v.Logger.Error(err.Error())
			}
		}

		if v.VlessManager != nil {
			if err := v.VlessManager.DelUser(user.Email); err != nil {
				v.Logger.Error(err.Error())
			}

		}

		if v.VmessManager != nil {
			if err := v.VmessManager.AddUser(user); err != nil {
				v.Logger.Error(err.Error())
			}
			v.Logger.Info("修改Vmess用户", zap.String("email", user.Email), zap.String("uuid", user.UUID))
		}

		if v.VlessManager != nil {
			if err := v.VlessManager.AddUser(user); err != nil {
				v.Logger.Error(err.Error())
			}
			v.Logger.Info("修改Vless用户", zap.String("email", user.Email), zap.String("uuid", user.UUID))
		}

	}

	for i := range removedUsers {
		ru := removedUsers[i]

		if v.VmessManager != nil {
			if err := v.VmessManager.DelUser(ru.Email); err != nil {
				v.Logger.Error(err.Error())
			}
			v.Logger.Info("删除Vmess用户", zap.String("email", ru.Email), zap.String("uuid", ru.UUID))
		}

		if v.VlessManager != nil {
			if err := v.VlessManager.DelUser(ru.Email); err != nil {
				v.Logger.Error(err.Error())
			}
			v.Logger.Info("删除Vless用户", zap.String("email", ru.Email), zap.String("uuid", ru.UUID))
		}

	}
	v.Accounts = newAccounts
}

func (v *VClient) syncTraffics(nodeId int64) error {
	var totalTraffic int64
	for _, a := range v.Accounts {
		ut, err := v.Stats.GetUserTraffic(a.Email, true)
		if err != nil {
			v.Logger.Error(err.Error())
			continue
		}
		a.Traffics = ut
		if ut.Uploads+ut.Downloads > 0 {
			if err := controllers.Agent.UpdateAccountTraffics(context.TODO(), nodeId, a); err != nil {
				v.Logger.Error(err.Error())
				continue
			}
			v.Logger.Debug("sync account traffic to database", zap.String("email", a.Email), zap.Int64("uplink", a.Traffics.Uploads), zap.Int64("downlink", a.Traffics.Downloads), zap.Int64("clients count", a.Traffics.Clients), zap.Strings("ips", a.Traffics.IPs))
			totalTraffic = totalTraffic + ut.Downloads + ut.Uploads
		}
	}
	if err := controllers.Agent.Heartbeat(context.TODO(), nodeId, totalTraffic); err != nil {
		v.Logger.Error(err.Error())
		return err
	}

	return nil
}
