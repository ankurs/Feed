package db

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound    = errors.New("Not found")
	ErrInvalidData = errors.New("Invalid Data")
)

const (
	FOLLOWING_FEED int32 = iota
	USER_FEED
)

type DB interface {
	AddFollowing(ctx context.Context, userId, followingId string) error
	AddFollower(ctx context.Context, userId, followerId string) error
	RemoveFollowing(ctx context.Context, userId, followingId string) error
	RemoveFollower(ctx context.Context, userId, followerId string) error
	CheckUserName(ctx context.Context, username string) (string, error)
	CheckEmail(ctx context.Context, email string) (string, error)
	CheckLogin(ctx context.Context, username, password string, hash func(context.Context, string, string) string) (UserInfo, error)
	CreateUser(ctx context.Context, req UserInfo, password string, hash func(context.Context, string, string) string) (string, error)
	GetUser(ctx context.Context, userID string) (UserInfo, error)
	CreateFeedItem(ctx context.Context, fi FeedInfo, ts time.Time) (string, error)
	AddUserFeedItem(ctx context.Context, userId, itemId string, ts time.Time) error
	AddFollowingFeedItem(ctx context.Context, userId, itemId string, ts time.Time) error
	GetFollowers(ctx context.Context, userId string) <-chan Data
	FetchFeed(ctx context.Context, userId string, before time.Time, ftype int32, limit int) ([]string, error)
	FetchFeedItem(ctx context.Context, feedId string) (FeedInfo, error)
	Close()
}

type Cache interface {
	GetUser(ctx context.Context, userID string) (UserInfo, error)
	GetFeedItem(ctx context.Context, feedId string) (FeedInfo, error)
	SetUser(ctx context.Context, ui UserInfo) error
	SetFeedItem(ctx context.Context, fi FeedInfo) error
}

type UserInfo interface {
	GetLastName() string
	GetFirstName() string
	GetUserName() string
	GetEmail() string
	GetId() string
}

type FeedInfo interface {
	GetId() string
	GetActor() string
	GetVerb() string
	GetCVerb() string
	GetObject() string
	GetTarget() string
	GetTs() int64
}

type Data interface {
	GetError() error
	GetId() string
}

func BuildInterface(vals ...interface{}) []interface{} {
	return vals
}
