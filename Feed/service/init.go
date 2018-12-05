package service

import (
	"github.com/carousell/Orion/orion"
	"github.com/mitchellh/mapstructure"
)

type svcFactory struct {
}

func (s *svcFactory) NewService(svr orion.Server) interface{} {
	cfg := Config{}
	if c, ok := svr.GetConfig()["feed"]; ok {
		mapstructure.Decode(c, &cfg)
	}
	return GetService(cfg)
}

func (s *svcFactory) DisposeService(svc interface{}) {
	DestroyService(svc)
}

func GetServiceFactory() orion.ServiceFactory {
	return &svcFactory{}
}

func RegisterOptionals(server orion.Server) {
}
