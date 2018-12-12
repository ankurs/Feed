package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ankurs/Feed/Feed/service/store/db"
	cproto "github.com/ankurs/Feed/Feed/service/store/redis/cache_proto"
	"github.com/carousell/Orion/utils/spanutils"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	userPrefix        string = "user::"
	feedPrefix        string = "feed::"
	usersFeedPrefix   string = "user:feed::"
	usersFollowPrefix string = "user:follow::"
)

type cache struct {
	client *redis.Client
	expiry time.Duration
}

func (r *cache) GetUser(ctx context.Context, userID string) (db.UserInfo, error) {
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

func (r *cache) GetFeedItem(ctx context.Context, feedId string) (db.FeedInfo, error) {
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

func (r *cache) SetUser(ctx context.Context, ui db.UserInfo) error {
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
		err = r.client.Set(key, data, r.expiry).Err()
		if err != nil {
			span.SetError(err.Error())
		}
		return errors.Wrap(err, "SetUser")
	}
	return db.ErrInvalidData
}

func (r *cache) SetFeedItem(ctx context.Context, fi db.FeedInfo) error {
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
		err = r.client.Set(key, data, r.expiry).Err()
		if err != nil {
			span.SetError(err.Error())
		}
		return errors.Wrap(err, "SetFeedItem")
	}
	return db.ErrInvalidData
}

func (r *cache) AddUserFeedItem(ctx context.Context, userId, itemId string, ts time.Time) error {
	if userId != "" && itemId != "" {
		name := "AddUserFeedItem"
		key := usersFeedPrefix + userId
		return r.addFeedItem(ctx, name, key, itemId, ts)
	}
	return db.ErrInvalidData
}

func (r *cache) AddFollowingFeedItem(ctx context.Context, userId, itemId string, ts time.Time) error {
	if userId != "" && itemId != "" {
		name := "AddFollowingFeedItem"
		key := usersFollowPrefix + userId
		return r.addFeedItem(ctx, name, key, itemId, ts)
	}
	return db.ErrInvalidData
}

func (r *cache) addFeedItem(ctx context.Context, name, key, itemId string, ts time.Time) error {
	span, ctx := spanutils.NewDatastoreSpan(ctx, name, "Redis")
	defer span.Finish()

	if key != "" {
		span.SetTag("key", key)
		z := redis.Z{
			Member: itemId,
			Score:  float64(ts.Unix()),
		}
		// TODO add size checks as well
		err := r.client.ZAdd(key, z).Err()
		if err != nil {
			span.SetError(err.Error())
			return errors.Wrap(err, name)
		}
	}
	return db.ErrInvalidData

}

func (r *cache) FetchFeed(ctx context.Context, userId string, before time.Time, ftype int32, limit int) ([]string, error) {
	span, ctx := spanutils.NewDatastoreSpan(ctx, "FetchFeed", "Redis")
	defer span.Finish()

	prefix := usersFollowPrefix
	if ftype == db.USER_FEED {
		prefix = usersFeedPrefix
	}

	if limit < 0 {
		limit = 20
	} else if limit > 50 {
		limit = 50
	}

	key := prefix + userId
	span.SetTag("key", key)
	opt := redis.ZRangeBy{
		Max:   strconv.FormatInt(before.Unix(), 10),
		Min:   "0",
		Count: int64(limit),
	}
	return r.client.ZRevRangeByScore(key, opt).Result()
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
	return &cache{
		client: redis.NewClient(&redis.Options{
			Addr: config.Host,
		}),
		expiry: time.Duration(config.Expiry) * time.Second,
	}
}
