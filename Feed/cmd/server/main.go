package main

import (
	proto "github.com/ankurs/Feed/Feed/Feed_proto"
	"github.com/ankurs/Feed/Feed/service"
	"github.com/carousell/Orion/orion"
)

func main() {
	server := orion.GetDefaultServer("EchoService")
	proto.RegisterFeedOrionServer(service.GetServiceFactory(), server)
	service.RegisterOptionals(server)
	server.Start()
	server.Wait()
}
