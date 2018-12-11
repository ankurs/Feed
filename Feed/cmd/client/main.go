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

	a := proto.NewAccountClient(conn)
	// make sure users are registered
	register1(a)
	register2(a)

	// get users from login
	u1, a1 := login(a, "ABCXYZ", "password")
	u2, _ := login(a, "LMNOPQ", "password")

	f := proto.NewFollowClient(conn)
	follow(f, u1, u2)
	follow(f, u2, u1)
	//unfollow(f, u1, u2)

	feed := proto.NewFeedClient(conn)
	addFeedItem(feed, u2)
	fetchFeed(feed, a1)
}

func login(c proto.AccountClient, user, pwd string) (*proto.UserInfo, *proto.Auth) {
	fmt.Println("making login gRPC call")
	req := new(proto.LoginRequest)
	req.UserName = user //"ABCXYZ"
	req.Password = pwd  //"password"
	r, err := c.Login(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Response : %s", r)
	return r.GetUser(), r.GetAuth()
}

func register1(c proto.AccountClient) *proto.RegisterResponse {
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
	return r
}

func register2(c proto.AccountClient) *proto.RegisterResponse {
	req := new(proto.RegisterRequest)
	req.FirstName = "LMN"
	req.LastName = "OPQ"
	req.UserName = "LMNOPQ"
	req.Email = "LMN@OPQ"
	req.Password = "password"
	r, err := c.Register(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Response : %s", r)
	return r
}

func follow(f proto.FollowClient, u1 *proto.UserInfo, u2 *proto.UserInfo) {
	fmt.Println(u1, u2)
	req := new(proto.FollowRequest)
	req.UserId = u1.GetId()
	req.FollowingId = u2.GetId()
	r, err := f.AddFollow(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Response : %s", r)
}

func unfollow(f proto.FollowClient, u1 *proto.UserInfo, u2 *proto.UserInfo) {
	req := new(proto.UnfollowRequest)
	req.UserId = u1.GetId()
	req.FollowingId = u2.GetId()
	r, err := f.RemoveFollow(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Response : %s", r)
}

func addFeedItem(f proto.FeedClient, user *proto.UserInfo) {
	req := new(proto.AddFeedItemRequest)
	req.Item = new(proto.FeedItem)
	req.Item.Actor = user.GetId()
	req.Item.Verb = proto.Verb_POST
	r, err := f.AddFeed(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Response : %s", r)
}

func fetchFeed(f proto.FeedClient, a *proto.Auth) {
	req := new(proto.FeedRequest)
	req.Auth = a
	req.Count = 10
	r, err := f.FetchFeed(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Response : %s", r)
}
