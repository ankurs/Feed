package store

import (
	"context"
	"strings"
	"time"

	"github.com/ankurs/Feed/Feed/service/store/cassandra"
	"github.com/ankurs/Feed/Feed/service/store/db"
	"github.com/carousell/Orion/utils/errors"
	"github.com/carousell/Orion/utils/log"
	"github.com/carousell/Orion/utils/spanutils"
)

type str struct {
	cas db.DB
}

var (
	ErrAlreadyTaken = errors.New("error already taken")
	ErrInvalidLogin = errors.New("could not login to the account")
	ErrNoSuchUser   = errors.New("no such user")
)

func (s *str) Register(ctx context.Context, req RegisterRequest) (LoginResponse, error) {
	name := "StorageRegister"
	// zipkin span
	span, ctx := spanutils.NewInternalSpan(ctx, name)
	defer span.Finish()

	username := strings.ToLower(req.GetUserName())
	_, err := s.cas.CheckUserName(ctx, username)
	if cause(err) != db.ErrNotFound {
		if err == nil {
			return nil, ErrAlreadyTaken
		}
		return nil, errors.Wrap(err, name)
	}

	_, err = s.cas.CheckEmail(ctx, req.GetEmail())
	if cause(err) != db.ErrNotFound {
		if err == nil {
			return nil, ErrAlreadyTaken
		}
		return nil, errors.Wrap(err, name)
	}
	ui := userInfo{
		firstname: req.GetFirstName(),
		lastname:  req.GetLastName(),
		email:     req.GetEmail(),
		username:  strings.ToLower(req.GetUserName()),
	}
	id, err := s.cas.CreateUser(ctx, ui, req.GetPassword(), getPasswordHash)
	if err != nil {
		return nil, errors.Wrap(err, name)
	}
	// update id in response
	ui.id = id
	resp := loginResponse{
		userInfo: ui,
	}
	return resp, nil
}

func (s *str) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	user, err := s.cas.CheckLogin(ctx, req.GetUserName(), req.GetPassword(), getPasswordHash)
	if err == nil && user != nil {
		// for now just use id as login token
		// TODO move to JWT token
		return loginResponse{token: user.GetId(), userInfo: user}, nil
	}
	if cause(err) == db.ErrNotFound {
		return nil, ErrInvalidLogin
	}
	return nil, errors.Wrap(err, "Login")
}

func (s *str) GetUser(ctx context.Context, userID string) (UserInfo, error) {
	ui, err := s.cas.GetUser(ctx, userID)
	if cause(err) == db.ErrNotFound {
		return nil, ErrNoSuchUser
	}
	return ui, errors.Wrap(err, "GetUser")
}

func (s *str) AddFollow(ctx context.Context, userId, followingId string) error {
	// TODO ensure this is transactional/add retry
	err := s.cas.AddFollowing(ctx, userId, followingId)
	if err != nil {
		return errors.Wrap(err, "AddFollowing")
	}
	err = s.cas.AddFollower(ctx, followingId, userId)
	if err != nil {
		return errors.Wrap(err, "AddFollower")
	}
	return nil
}

func (s *str) RemoveFollow(ctx context.Context, userId, followingId string) error {
	// TODO ensure this is transactional/add retry
	err := s.cas.RemoveFollowing(ctx, userId, followingId)
	if err != nil {
		return errors.Wrap(err, "RemoveFollowing")
	}
	err = s.cas.RemoveFollower(ctx, followingId, userId)
	if err != nil {
		return errors.Wrap(err, "RemoveFollower")
	}
	return nil
}

func (s *str) CreateFeedItem(ctx context.Context, fi FeedInfo, ts time.Time) (string, error) {
	id, err := s.cas.CreateFeedItem(ctx, fi, ts)
	if err != nil {
		return "", errors.Wrap(err, "AddFeedItem")
	}
	return id, nil
}

func (s *str) AddUserFeedItem(ctx context.Context, userId, itemId string, ts time.Time) error {
	err := s.cas.AddUserFeedItem(ctx, userId, itemId, ts)
	if err != nil {
		return errors.Wrap(err, "AddUserFeedItem")
	}
	return nil
}

func (s *str) GetFollowers(ctx context.Context, userId string) ([]string, error) {
	// TODO implement this as an iterator
	followers, err := s.cas.GetFollowers(ctx, userId)
	if err != nil {
		return followers, errors.Wrap(err, "GetFollowers")
	}
	return followers, nil
}

func (s *str) AddFollowingFeedItem(ctx context.Context, userId, itemId string, ts time.Time) error {
	err := s.cas.AddFollowingFeedItem(ctx, userId, itemId, ts)
	if err != nil {
		return errors.Wrap(err, "AddUserFeedItem")
	}
	return nil
}

func (s *str) Close() {
	if s != nil && s.cas != nil {
		s.cas.Close()
	}
}

func NewStore(config Config) (Storage, error) {
	cas, err := cassandra.New(config.Cassandra)
	if err != nil {
		log.Error(context.Background(), err)
		return nil, errors.Wrap(err, "NewStore")
	}
	return &str{
		cas: cas,
	}, nil
}
