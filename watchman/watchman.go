package watchman

import (
	"time"

	"github.com/xtls/xray-core/watchman/logging"
	"github.com/xtls/xray-core/watchman/vclient"
)

var logger = logging.GetInstance().Logger.Sugar()

type Server struct {
	Config *Config
}

func (w *Server) Start() {
	logger.Debugf("watchman start, %+v", *w.Config)
	vc, err := vclient.Connect(w.Config.ApiAddress, time.Second*6)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("watchman connected to v2ray")
	if err := vc.InitServices(w.Config.VmessInboundTag, w.Config.VlessInboundTag); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("watchman services start")
	if err := vc.Startup(w.Config.DBUrl, w.Config.NodeID, w.Config.CheckRate, w.Config.VIP); err != nil {
		logger.Fatal(err.Error())
	}
}

func New(config *Config) *Server {
	return &Server{
		Config: config,
	}
}
