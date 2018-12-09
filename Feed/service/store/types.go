package store

import (
	"context"
	"time"
)

type RegisterRequest interface {
	GetLastName() string
	GetFirstName() string
	GetUserName() string
	GetPassword() string
	GetEmail() string
}

type RegisterResponse interface {
	GetId() string
}

type LoginRequest interface {
	GetUserName() string
	GetPassword() string
}

type LoginResponse interface {
	Gettoken() string
}

type Storage interface {
	Register(context.Context, RegisterRequest) (RegisterResponse, error)
	Login(context.Context, LoginRequest) (LoginResponse, error)
	Close()
}

type Config struct {
	//CassandraHosts are the hosts that storage will connect to
	CassandraHosts []string
	//CassandraConsistency is the consistency level for all cassandra calls (ideally this should be set to 'LOCAL_QUORUM')
	CassandraConsistency string
	//CassandraConnectTimeout is the time initial connection to cassandra will wait before timing out
	CassandraConnectTimeout time.Duration
	//CassandraOperationTimeout is the time each operation will wait before timing out
	CassandraOperationTimeout time.Duration
	//NumConns is the number of connections that are maintained per cassandra hosts
	NumConns int
}
