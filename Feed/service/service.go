// Package service must implement the generated proto's server interface
package service

import (
	"context"
	"time"

	proto "github.com/ankurs/Feed/Feed/Feed_proto"
	"github.com/ankurs/Feed/Feed/service/store"
	"github.com/carousell/Orion/utils/errors"
	"github.com/carousell/Orion/utils/log"
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
	return statusCustom(200, false, "")
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
		resp.User = new(proto.UserInfo)
		resp.User.UserName = r.GetUserInfo().GetUserName()
		resp.User.FirstName = r.GetUserInfo().GetFirstName()
		resp.User.LastName = r.GetUserInfo().GetLastName()
		resp.User.Email = r.GetUserInfo().GetEmail()
		resp.User.Id = r.GetUserInfo().GetId()
	}
	return resp, nil
}

func (s *svc) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	resp := new(proto.LoginResponse)
	r, err := s.storage.Login(ctx, req)
	if err != nil && err != store.ErrInvalidLogin {
		return nil, errors.WrapWithStatus(err, "Login", status.New(codes.Internal, "error logging in"))
	}
	if err == store.ErrInvalidLogin {
		resp.Status = statusCustom(401, true, "invalid username/password")
		return resp, nil
	}
	resp.Auth = new(proto.Auth)
	resp.Auth.Token = r.GetToken()
	resp.User = new(proto.UserInfo)
	resp.User.UserName = r.GetUserInfo().GetUserName()
	resp.User.FirstName = r.GetUserInfo().GetFirstName()
	resp.User.LastName = r.GetUserInfo().GetLastName()
	resp.User.Email = r.GetUserInfo().GetEmail()
	resp.User.Id = r.GetUserInfo().GetId()
	return resp, nil
}

func (s *svc) FetchFeed(context.Context, *proto.FeedRequest) (*proto.FeedResponse, error) {
	panic("not implemented")
}

func (s *svc) AddFeed(ctx context.Context, req *proto.AddFeedItemRequest) (*proto.AddFeedItemResponse, error) {
	resp := new(proto.AddFeedItemResponse)
	ts := time.Now()
	fi := feedInfo{*req.GetItem()}
	// TODO validate user id
	id, err := s.storage.CreateFeedItem(ctx, &fi, ts)
	if err != nil {
		resp.Status = statusCustom(500, true, "could not add feed: "+err.Error())
		return resp, nil
	}
	err = s.addFeedToUser(ctx, fi.GetActor(), id, ts)
	if err != nil {
		resp.Status = statusCustom(500, true, "could not add feed: "+err.Error())
		return resp, nil
	}
	go s.addFeedItemToFollowes(ctx, fi.GetActor(), id, ts)
	resp.Status = statusOK()
	resp.Id = id
	return resp, nil
}

func (s *svc) addFeedToUser(ctx context.Context, userId, feedId string, ts time.Time) error {
	return s.storage.AddUserFeedItem(ctx, userId, feedId, ts)
}

// TODO move to worker
func (s *svc) addFeedItemToFollowes(ctx context.Context, userId, feedId string, ts time.Time) {
	followers, err := s.storage.GetFollowers(ctx, userId)
	if err != nil {
		//TODO handle failure with some retry
		log.Error(ctx, "err", err)
		return
	}
	for _, follower := range followers {
		//TODO handle failure with some retry
		s.storage.AddFollowingFeedItem(ctx, follower, feedId, ts)
	}
}

func (s *svc) AddFollow(ctx context.Context, req *proto.FollowRequest) (*proto.FollowResponse, error) {
	resp := new(proto.FollowResponse)
	// TODO validate user id
	err := s.storage.AddFollow(ctx, req.GetUserId(), req.GetFollowingId())
	if err != nil {
		resp.Status = statusCustom(500, true, "could not follow: "+err.Error())
	} else {
		resp.Status = statusOK()
	}
	return resp, nil
}

func (s *svc) RemoveFollow(ctx context.Context, req *proto.UnfollowRequest) (*proto.UnfollowResponse, error) {
	resp := new(proto.UnfollowResponse)
	// TODO validate user id
	err := s.storage.RemoveFollow(ctx, req.GetUserId(), req.GetFollowingId())
	if err != nil {
		resp.Status = statusCustom(500, true, "could not unfollow: "+err.Error())
	} else {
		resp.Status = statusOK()
	}
	return resp, nil
}

func (s *svc) Close() {
	// do nothing
}

func DestroyService(obj interface{}) {
	if s, ok := obj.(interface{ Close() }); ok {
		s.Close()
	}
}
