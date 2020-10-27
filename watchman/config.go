package watchman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/infra/conf"
	json_reader "v2ray.com/core/infra/conf/json"
	"v2ray.com/core/main/confloader"
)

const API_ADDRESS string = "127.0.0.1:4321"

type Config struct {
	NodeID          int64  `json:"nodeId"`
	CheckRate       int64  `json:"checkRate"`
	DBUrl           string `json:"dburl"`
	ApiAddress      string `json:"apiaddress"`
	VmessInboundTag string `json:"vmessInboundTag"`
	VlessInboundTag string `json:"vlessInboundTag"`
	v2rayConfig     *conf.Config
}

func LoadConfig(configFile string) (*Config, error) {
	type config struct {
		*conf.Config
		Watchman *Config `json:"watchman"`
	}

	configInput, err := confloader.LoadConfig(configFile)
	if err != nil {
		return nil, errors.New("failed to load config: ", configFile).Base(err)
	}

	cfg := &config{}
	if err = decodeCommentJSON(configInput, cfg); err != nil {
		return nil, err
	}
	if cfg.Watchman != nil {
		cfg.Watchman.v2rayConfig = cfg.Config
		if err = checkCfg(cfg.Watchman); err != nil {
			return nil, err
		}

		if cfg.Watchman.ApiAddress == "" {
			cfg.Watchman.ApiAddress = API_ADDRESS
		}
	}

	return cfg.Watchman, err
}

func checkCfg(cfg *Config) error {

	if cfg.v2rayConfig.Api == nil {
		return errors.New("Api must be set")
	}

	apiTag := cfg.v2rayConfig.Api.Tag
	if len(apiTag) == 0 {
		return errors.New("Api tag can't be empty")
	}

	services := cfg.v2rayConfig.Api.Services
	if !InStr("HandlerService", services) {
		return errors.New("Api service, HandlerService, must be enabled")
	}
	if !InStr("StatsService", services) {
		return errors.New("Api service, StatsService, must be enabled")
	}

	if cfg.v2rayConfig.Stats == nil {
		return errors.New("Stats must be enabled")
	}

	if apiInbound := getInboundConfigByTag(apiTag, cfg.v2rayConfig.InboundConfigs); apiInbound == nil {
		return errors.New(fmt.Sprintf("Miss an inbound tagged %s", apiTag))
	} else if apiInbound.Protocol != "dokodemo-door" {
		return errors.New(fmt.Sprintf("The protocol of inbound tagged %s must be \"dokodemo-door\"", apiTag))
	} else {
		if apiInbound.ListenOn == nil || apiInbound.PortRange == nil {
			return errors.New(fmt.Sprintf("Fields, \"listen\" and \"port\", of inbound tagged %s must be set", apiTag))
		}
	}

	if inbound := getInboundConfigByTag(cfg.VmessInboundTag, cfg.v2rayConfig.InboundConfigs); inbound == nil {
		return errors.New(fmt.Sprintf("Miss an inbound tagged %s", cfg.VmessInboundTag))
	}

	if inbound := getInboundConfigByTag(cfg.VlessInboundTag, cfg.v2rayConfig.InboundConfigs); inbound == nil {
		return errors.New(fmt.Sprintf("Miss an inbound tagged %s", cfg.VlessInboundTag))
	}

	return nil
}

func getInboundConfigByTag(apiTag string, inbounds []conf.InboundDetourConfig) *conf.InboundDetourConfig {
	for _, inbound := range inbounds {
		if inbound.Tag == apiTag {
			return &inbound
		}
	}
	return nil
}

func decodeCommentJSON(reader io.Reader, i interface{}) error {
	jsonContent := bytes.NewBuffer(make([]byte, 0, 10240))
	jsonReader := io.TeeReader(&json_reader.Reader{
		Reader: reader,
	}, jsonContent)
	decoder := json.NewDecoder(jsonReader)
	return decoder.Decode(i)
}

func InStr(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
