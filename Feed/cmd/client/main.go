package main

import (
	"context"
	"fmt"
	"log"

	proto "github.com/ankurs/Feed/Feed/Feed_proto"
	"google.golang.org/grpc"
)

const (
	address = "192.168.99.100:9281"
	//address = "127.0.0.1:9281"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewAccountClient(conn)
	register(c)
	login(c)
}

func login(c proto.AccountClient) {
	fmt.Println("making login gRPC call")
	req := new(proto.LoginRequest)
	req.UserName = "ABCXYZ"
	req.Password = "password"
	r, err := c.Login(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Response : %s", r)
}

func register(c proto.AccountClient) {
	req := new(proto.RegisterRequest)
	req.FirstName = "ABC"
	req.LastName = "XYZ"
	req.UserName = "ABCXYZ"
	req.Email = "ABC@XYZ"
	req.Password = "password"
	r, err := c.Register(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Response : %s", r)
}
