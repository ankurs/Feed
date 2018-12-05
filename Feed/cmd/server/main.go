package main

import (
	proto "github.com/ankurs/Feed/Feed/Feed_proto"
	"github.com/ankurs/Feed/Feed/service"
	"github.com/carousell/Orion/orion"
	"github.com/carousell/Orion/orion/helpers"
)

func main() {
	server := orion.GetDefaultServer("feed")

	factory, err := helpers.NewSingleServiceFactory(service.GetServiceFactory())
	if err != nil {
		panic(err)
	}

	// register services
	proto.RegisterFeedOrionServer(factory, server)
	proto.RegisterAccountOrionServer(factory, server)

	// register optionals
	service.RegisterOptionals(server)

	server.Start()
	server.Wait()
}
