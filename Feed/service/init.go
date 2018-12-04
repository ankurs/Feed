package service

import (
	"fmt"

	proto "github.com/ankurs/Feed/Feed/Feed_proto"
	"github.com/carousell/Orion/orion"
	"github.com/mitchellh/mapstructure"
)

type svcFactory struct {
}

func (s *svcFactory) NewService(svr orion.Server) interface{} {
	cfg := Config{}
	if c, ok := svr.GetConfig()["echo"]; ok {
		mapstructure.Decode(c, &cfg)
	}
	return GetService(cfg)
}

func (s *svcFactory) DisposeService(svc interface{}) {
	fmt.Println("disposing", svc)
	DestroyService(svc)
}

func GetServiceFactory() orion.ServiceFactory {
	return &svcFactory{}
}

func RegisterOptionals(server orion.Server) {
	proto.RegisterFeedUpperEncoder(server, encoder)
	proto.RegisterFeedUpperDecoder(server, decoder)
	proto.RegisterFeedUpperHandler(server, optionsHandler)
}
