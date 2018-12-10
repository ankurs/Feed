package store

import (
	"context"

	"github.com/ankurs/Feed/Feed/service/store/cassandra"
	"github.com/ankurs/Feed/Feed/service/store/db"
)

type RegisterRequest interface {
	GetLastName() string
	GetFirstName() string
	GetUserName() string
	GetPassword() string
	GetEmail() string
}

type LoginRequest interface {
	GetUserName() string
	GetPassword() string
}

type LoginResponse interface {
	GetToken() string
	GetUserInfo() UserInfo
}

// we type alias it, so that we can saperate them out in future
type UserInfo = db.UserInfo

type Storage interface {
	Register(context.Context, RegisterRequest) (LoginResponse, error)
	Login(context.Context, LoginRequest) (LoginResponse, error)
	GetUser(ctx context.Context, userID string) (UserInfo, error)
	AddFollow(ctx context.Context, userId, followingId string) error
	RemoveFollow(ctx context.Context, userId, followingId string) error
	Close()
}

type Config struct {
	Cassandra cassandra.Config
}
