// Package service must implement the generated proto's server interface
package service

import (
	"context"

	proto "github.com/ankurs/Feed/Feed/Feed_proto"
	"github.com/carousell/Orion/utils/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type svc struct {
}

func (s *svc) Register(context.Context, *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return nil, errors.NewWithStatus("not implemented yet", status.New(codes.Unimplemented, "not implemented yet"))
}

func (s *svc) Login(context.Context, *proto.LoginRequest) (*proto.LoginResponse, error) {
	return nil, errors.NewWithStatus("not implemented yet", status.New(codes.Unimplemented, "not implemented yet"))
}

func (s *svc) Fetch(context.Context, *proto.FeedRequest) (*proto.FeedResponse, error) {
	return nil, errors.NewWithStatus("not implemented yet", status.New(codes.Unimplemented, "not implemented yet"))
}

func (s *svc) Close() {
	// do nothing
}

func GetService(config Config) proto.FeedServer {
	s := new(svc)
	return s
}

func DestroyService(obj interface{}) {
	if s, ok := obj.(interface{ Close() }); ok {
		s.Close()
	}
}
