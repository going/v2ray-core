package vclient

import (
	"context"

	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/common/uuid"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/shadowsocks"
	"github.com/xtls/xray-core/proxy/vless"
	"github.com/xtls/xray-core/proxy/vmess"
	"github.com/xtls/xray-core/proxy/vmess/inbound"
	"github.com/xtls/xray-core/transport/internet"
	"github.com/xtls/xray-core/transport/internet/headers/noop"
	"github.com/xtls/xray-core/transport/internet/headers/srtp"
	"github.com/xtls/xray-core/transport/internet/headers/tls"
	"github.com/xtls/xray-core/transport/internet/headers/utp"
	"github.com/xtls/xray-core/transport/internet/headers/wechat"
	"github.com/xtls/xray-core/transport/internet/headers/wireguard"
	"github.com/xtls/xray-core/transport/internet/kcp"
	"github.com/xtls/xray-core/transport/internet/websocket"
	"github.com/xtls/xray-core/watchman/proto"
	"google.golang.org/grpc"
)

var KcpHeadMap = map[string]*serial.TypedMessage{
	"wechat-video": serial.ToTypedMessage(&wechat.VideoConfig{}),
	"srtp":         serial.ToTypedMessage(&srtp.Config{}),
	"utp":          serial.ToTypedMessage(&utp.Config{}),
	"wireguard":    serial.ToTypedMessage(&wireguard.WireguardConfig{}),
	"dtls":         serial.ToTypedMessage(&tls.PacketConfig{}),
	"noop":         serial.ToTypedMessage(&noop.Config{}),
}
var CipherTypeMap = map[string]shadowsocks.CipherType{
	"aes-256-cfb":            shadowsocks.CipherType_AES_256_CFB,
	"aes-128-cfb":            shadowsocks.CipherType_AES_128_CFB,
	"aes-128-gcm":            shadowsocks.CipherType_AES_128_GCM,
	"aes-256-gcm":            shadowsocks.CipherType_AES_256_GCM,
	"chacha20":               shadowsocks.CipherType_CHACHA20,
	"chacah-ietf":            shadowsocks.CipherType_CHACHA20_IETF,
	"chacha20-ploy1305":      shadowsocks.CipherType_CHACHA20_POLY1305,
	"chacha20-ietf-poly1305": shadowsocks.CipherType_CHACHA20_POLY1305,
}

type HandlerServiceClient struct {
	command.HandlerServiceClient
	InboundTag   string
	VlessEnabled bool
}

func NewHandlerServiceClient(client *grpc.ClientConn, inboundTag string, vlessEnabled bool) *HandlerServiceClient {
	return &HandlerServiceClient{
		HandlerServiceClient: command.NewHandlerServiceClient(client),
		InboundTag:           inboundTag,
		VlessEnabled:         vlessEnabled,
	}
}

// user
func (h *HandlerServiceClient) DelUser(email string) error {
	req := &command.AlterInboundRequest{
		Tag:       h.InboundTag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{Email: email}),
	}
	return h.AlterInbound(req)
}

func (h *HandlerServiceClient) AddUser(user *proto.UserModel) error {
	req := &command.AlterInboundRequest{
		Tag:       h.InboundTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{User: h.ConvertVmessUser(user)}),
	}
	if h.VlessEnabled {
		req.Operation = serial.ToTypedMessage(&command.AddUserOperation{User: h.ConvertVlessUser(user)})
	}

	return h.AlterInbound(req)
}

func (h *HandlerServiceClient) AlterInbound(req *command.AlterInboundRequest) error {
	_, err := h.HandlerServiceClient.AlterInbound(context.Background(), req)
	return err
}

//streaming
func GetKcpStreamConfig(headkey string) *internet.StreamConfig {
	var streamsetting internet.StreamConfig
	head, _ := KcpHeadMap["noop"]
	if _, ok := KcpHeadMap[headkey]; ok {
		head, _ = KcpHeadMap[headkey]
	}
	streamsetting = internet.StreamConfig{
		ProtocolName: "mkcp",
		TransportSettings: []*internet.TransportConfig{
			{
				ProtocolName: "mkcp",
				Settings: serial.ToTypedMessage(
					&kcp.Config{
						HeaderConfig: head,
					}),
			},
		},
	}
	return &streamsetting
}

func GetWebSocketStreamConfig(path string, host string) *internet.StreamConfig {
	var streamsetting internet.StreamConfig
	streamsetting = internet.StreamConfig{
		ProtocolName: "websocket",
		TransportSettings: []*internet.TransportConfig{
			{
				ProtocolName: "websocket",
				Settings: serial.ToTypedMessage(&websocket.Config{
					Path: path,
					Header: []*websocket.Header{
						{
							Key:   "Hosts",
							Value: host,
						},
					},
				}),
			},
		},
	}
	return &streamsetting
}

// different type inbounds
func (h *HandlerServiceClient) AddVmessInbound(port uint16, address string, streamsetting *internet.StreamConfig) error {
	addInboundRequest := &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			Tag: h.InboundTag,
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortRange:      net.SinglePortRange(net.Port(port)),
				Listen:         net.NewIPOrDomain(net.ParseAddress(address)),
				StreamSettings: streamsetting,
			}),
			ProxySettings: serial.ToTypedMessage(&inbound.Config{
				User: []*protocol.User{
					{
						Level: 0,
						Email: "admin@tian.network",
						Account: serial.ToTypedMessage(&vmess.Account{
							Id:      protocol.NewID(uuid.New()).String(),
							AlterId: 2,
						}),
					},
				},
			}),
		},
	}
	return h.AddInbound(addInboundRequest)
}

func (h *HandlerServiceClient) AddInbound(req *command.AddInboundRequest) error {
	_, err := h.HandlerServiceClient.AddInbound(context.Background(), req)
	return err
}
func (h *HandlerServiceClient) RemoveInbound() error {
	req := &command.RemoveInboundRequest{
		Tag: h.InboundTag,
	}
	_, err := h.HandlerServiceClient.RemoveInbound(context.Background(), req)
	return err
}

func (h *HandlerServiceClient) ConvertVmessUser(userModel *proto.UserModel) *protocol.User {
	return &protocol.User{
		Level: 0,
		Email: userModel.Email,
		Account: serial.ToTypedMessage(&vmess.Account{
			Id:      userModel.UUID,
			AlterId: userModel.AlterID,
			SecuritySettings: &protocol.SecurityConfig{
				Type: protocol.SecurityType_AUTO,
			},
		}),
	}
}

func (h *HandlerServiceClient) ConvertVlessUser(userModel *proto.UserModel) *protocol.User {
	return &protocol.User{
		Level: 0,
		Email: userModel.Email,
		Account: serial.ToTypedMessage(&vless.Account{
			Id: userModel.UUID,
		}),
	}
}
