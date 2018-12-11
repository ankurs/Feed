package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/ankurs/Feed/Feed/service/store/db"
	cproto "github.com/ankurs/Feed/Feed/service/store/redis/cache_proto"
	"github.com/carousell/Orion/utils/spanutils"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	userPrefix string = "user::"
	feedPrefix string = "feed::"
)

type r struct {
	client *redis.Client
	expiry time.Duration
}

func (r *r) GetUser(ctx context.Context, userID string) (db.UserInfo, error) {
	span, ctx := spanutils.NewDatastoreSpan(ctx, "GetUser", "Redis")
	defer span.Finish()
	key := userPrefix + userID
	span.SetTag("key", key)
	result := r.client.Get(key)
	err := result.Err()
	if err != nil {
		span.SetError(err.Error())
		if result.Err() == redis.Nil {
			return nil, db.ErrNotFound
		}
		return nil, errors.Wrap(err, "GetUser")
	}
	data, _ := result.Bytes()
	userInfo := new(cproto.UserInfo)
	err = proto.Unmarshal(data, userInfo)
	if err != nil {
		span.SetError(err.Error())
		return nil, errors.Wrap(err, "GetUser")
	}
	return userInfo, nil
}

func (r *r) GetFeedItem(ctx context.Context, feedId string) (db.FeedInfo, error) {
	span, ctx := spanutils.NewDatastoreSpan(ctx, "GetFeedItem", "Redis")
	defer span.Finish()
	key := feedPrefix + feedId
	span.SetTag("key", key)
	result := r.client.Get(key)
	err := result.Err()
	if err != nil {
		span.SetError(err.Error())
		if result.Err() == redis.Nil {
			return nil, db.ErrNotFound
		}
		return nil, errors.Wrap(err, "GetFeedItem")
	}
	data, _ := result.Bytes()
	feedInfo := new(cproto.FeedInfo)
	err = proto.Unmarshal(data, feedInfo)
	if err != nil {
		span.SetError(err.Error())
		return nil, errors.Wrap(err, "GetFeedItem")
	}
	return feedInfo, nil
}

func (r *r) SetUser(ctx context.Context, ui db.UserInfo) error {
	span, ctx := spanutils.NewDatastoreSpan(ctx, "SetUser", "Redis")
	defer span.Finish()

	if ui != nil && ui.GetId() != "" {
		key := userPrefix + ui.GetId()
		span.SetTag("key", key)
		userInfo := &cproto.UserInfo{
			Id:        ui.GetId(),
			FirstName: ui.GetFirstName(),
			LastName:  ui.GetLastName(),
			Email:     ui.GetEmail(),
			UserName:  ui.GetUserName(),
		}
		data, err := proto.Marshal(userInfo)
		if err != nil {
			span.SetError(err.Error())
			return errors.Wrap(err, "SetUser")
		}
		err = r.client.SetNX(key, data, r.expiry).Err()
		if err != nil {
			span.SetError(err.Error())
		}
		return errors.Wrap(err, "SetUser")
	}
	return db.ErrInvalidData
}

func (r *r) SetFeedItem(ctx context.Context, fi db.FeedInfo) error {
	span, ctx := spanutils.NewDatastoreSpan(ctx, "SetFeedItem", "Redis")
	defer span.Finish()

	if fi != nil && fi.GetId() != "" {
		key := feedPrefix + fi.GetId()
		span.SetTag("key", key)
		feedInfo := &cproto.FeedInfo{
			Id:     fi.GetId(),
			Actor:  fi.GetActor(),
			Verb:   fi.GetVerb(),
			CVerb:  fi.GetCVerb(),
			Object: fi.GetObject(),
			Target: fi.GetTarget(),
			Ts:     fi.GetTs(),
		}
		data, err := proto.Marshal(feedInfo)
		if err != nil {
			span.SetError(err.Error())
			return errors.Wrap(err, "SetFeedItem")
		}
		err = r.client.SetNX(key, data, r.expiry).Err()
		if err != nil {
			span.SetError(err.Error())
		}
		return errors.Wrap(err, "SetFeedItem")
	}
	return db.ErrInvalidData
}

type Config struct {
	Host   string
	Expiry int64
}

func NewRedisCache(config Config) db.Cache {
	fmt.Println(config)
	if config.Expiry < 1 {
		config.Expiry = 5 * 60 // 5 min
	}
	return &r{
		client: redis.NewClient(&redis.Options{
			Addr: config.Host,
		}),
		expiry: time.Duration(config.Expiry) * time.Second,
	}
}
