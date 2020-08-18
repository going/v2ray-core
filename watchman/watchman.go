package watchman

import (
	"time"

	"v2ray.com/core/watchman/logging"
	"v2ray.com/core/watchman/vclient"
)

var logger = logging.GetInstance().Logger.Sugar()

func Start(address, vmessInboundTag, vlessInboundTag, dbUrl string, nodeId int64) {
	logger.Debug("watchman start")
	vc, err := vclient.Connect(address, time.Second*6)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("watchman connected to v2ray")
	if err := vc.InitServices(vmessInboundTag, vlessInboundTag); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("watchman services start")
	if err := vc.Startup(dbUrl, nodeId); err != nil {
		logger.Fatal(err.Error())
	}
}
