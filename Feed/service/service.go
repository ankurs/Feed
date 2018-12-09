// Package service must implement the generated proto's server interface
package service

import (
	"context"

	proto "github.com/ankurs/Feed/Feed/Feed_proto"
	"github.com/ankurs/Feed/Feed/service/store"
	"github.com/carousell/Orion/utils/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type svc struct {
	config  Config
	storage store.Storage
}

func GetService(config Config) FeedService {
	str, _ := store.NewStore(config.Store)
	s := new(svc)
	s.config = config
	s.storage = str
	return s
}

func getTestService() FeedService {
	return new(svc)
}

func statusOK() *proto.StatusResponse {
	return statusCustom(0, false, "")
}

func statusCustom(code int32, err bool, msg string) *proto.StatusResponse {
	resp := new(proto.StatusResponse)
	resp.Code = code
	resp.Error = err
	resp.Msg = msg
	return resp
}

func (s *svc) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	r, err := s.storage.Register(ctx, req)
	if err != nil && err != store.ErrAlreadyTaken {
		return nil, errors.WrapWithStatus(err, "Register", status.New(codes.Internal, "error registering user"))
	}
	resp := new(proto.RegisterResponse)
	if err == store.ErrAlreadyTaken {
		resp.Status = statusCustom(409, true, "username/email already present")
	} else {
		resp.Status = statusOK()
		resp.Id = r.GetId()
	}
	return resp, nil
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

func DestroyService(obj interface{}) {
	if s, ok := obj.(interface{ Close() }); ok {
		s.Close()
	}
}
