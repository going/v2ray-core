package vclient

import (
	"context"
	"fmt"
	"strings"

	"v2ray.com/core/watchman/proto"

	"google.golang.org/grpc"
	"v2ray.com/core/app/stats/command"
)

const (
	DownloadTraffic = ">>>downlink"
	UploadTraffic   = ">>>uplink"
)

type StatsServiceClient struct {
	command.StatsServiceClient
}

func NewStatsServiceClient(client *grpc.ClientConn) *StatsServiceClient {
	return &StatsServiceClient{
		StatsServiceClient: command.NewStatsServiceClient(client),
	}
}

func (s *StatsServiceClient) GetUserTraffic(email string, reset bool) (*proto.UserTrafficLog, error) {
	req := &command.QueryStatsRequest{
		Pattern: fmt.Sprintf("user>>>%s>>>traffic", email),
		Reset_:  reset,
	}

	resp, err := s.QueryStats(context.Background(), req)
	if err != nil {
		return nil, err
	}

	traffic := &proto.UserTrafficLog{
		Email: email,
	}

	for _, stat := range resp.Stat {
		if strings.Contains(stat.Name, DownloadTraffic) && stat.Value > 0 {
			traffic.Downloads += stat.Value
		}
		if strings.Contains(stat.Name, UploadTraffic) && stat.Value > 0 {
			traffic.Uploads += stat.Value
			// name := "user>>>" + user.Email + ">>>traffic>>>" + sessionInbound.Source.Address.String() + ">>>uplink"
			fields := strings.Split(stat.Name, ">>>")
			if len(fields) == 5 { // 获取用户登录IP
				if !contains(traffic.IPs, fields[3]) {
					traffic.IPs = append(traffic.IPs, fields[3])
					traffic.Clients += 1
				}
			}
		}
	}

	return traffic, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
