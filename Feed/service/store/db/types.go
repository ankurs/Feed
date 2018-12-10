package db

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("Not found")
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
	Close()
}

type UserInfo interface {
	GetLastName() string
	GetFirstName() string
	GetUserName() string
	GetEmail() string
	GetId() string
}

func BuildInterface(vals ...interface{}) []interface{} {
	return vals
}
