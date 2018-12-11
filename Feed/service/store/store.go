package store

import (
	"context"
	"strings"
	"time"

	"github.com/ankurs/Feed/Feed/service/store/cassandra"
	"github.com/ankurs/Feed/Feed/service/store/db"
	"github.com/ankurs/Feed/Feed/service/store/redis"
	"github.com/carousell/Orion/utils/errors"
	"github.com/carousell/Orion/utils/log"
	"github.com/carousell/Orion/utils/spanutils"
)

type str struct {
	cas   db.DB
	cache db.Cache
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
	s.cache.SetUser(ctx, ui)
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
	// fetch from cache
	ui, err := s.cache.GetUser(ctx, userID)
	if err == nil {
		return ui, nil
	}
	ui, err = s.cas.GetUser(ctx, userID)
	if cause(err) == db.ErrNotFound {
		return nil, ErrNoSuchUser
	}
	if err != nil {
		// set it in cache
		s.cache.SetUser(ctx, ui)
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
	s.cache.SetFeedItem(ctx, &feedIder{FeedInfo: fi, id: id})
	return id, nil
}

func (s *str) AddUserFeedItem(ctx context.Context, userId, itemId string, ts time.Time) error {
	err := s.cas.AddUserFeedItem(ctx, userId, itemId, ts)
	if err != nil {
		return errors.Wrap(err, "AddUserFeedItem")
	}
	return nil
}

func (s *str) GetFollowers(ctx context.Context, userId string) <-chan db.Data {
	return s.cas.GetFollowers(ctx, userId)
}

func (s *str) AddFollowingFeedItem(ctx context.Context, userId, itemId string, ts time.Time) error {
	err := s.cas.AddFollowingFeedItem(ctx, userId, itemId, ts)
	if err != nil {
		return errors.Wrap(err, "AddUserFeedItem")
	}
	return nil
}

func (s *str) FetchFeed(ctx context.Context, userId string, before time.Time, ftype int32, limit int) ([]FeedInfo, error) {
	feedInfo := make([]FeedInfo, 0)
	feeds, err := s.cas.FetchFeed(ctx, userId, before, ftype, limit)
	if err != nil {
		return feedInfo, errors.Wrap(err, "FetchFeed")
	}
	for _, feed := range feeds {
		fi, err := s.cache.GetFeedItem(ctx, feed)
		if err != nil {
			fi, err = s.cas.FetchFeedItem(ctx, feed)
			if err != nil {
				return []FeedInfo{}, errors.Wrap(err, "FetchFeed")
			}
			s.cache.SetFeedItem(ctx, fi)
		}
		feedInfo = append(feedInfo, fi)
	}
	return feedInfo, nil
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
		cas:   cas,
		cache: redis.NewRedisCache(config.Redis),
	}, nil
}
