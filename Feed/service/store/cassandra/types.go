package cassandra

import (
	"context"
	"time"

	"github.com/ankurs/Feed/Feed/service/store/db"
	"github.com/gocql/gocql"
)

type Cassandra interface {
	db.DB
	CassandraExec(ctx context.Context, name, query string, values []interface{}, dest []interface{}) error
	CassandraExecWithConsistency(ctx context.Context, name, query string, values []interface{}, dest []interface{}, cons gocql.Consistency) error
}

type Config struct {
	//Hosts are the hosts that storage will connect to
	Hosts []string
	//Consistency is the consistency level for all cassandra calls (ideally this should be set to 'LOCAL_QUORUM')
	Consistency string
	//ConnectTimeout is the time initial connection to cassandra will wait before timing out
	ConnectTimeout time.Duration
	//OperationTimeout is the time each operation will wait before timing out
	OperationTimeout time.Duration
	//NumConns is the number of connections that are maintained per cassandra hosts
	NumConns int
}

// implements db.UserInfo
type userInfo struct {
	lastname  string
	firstname string
	username  string
	email     string
	id        string
}

func (u userInfo) GetLastName() string {
	return u.lastname
}

func (u userInfo) GetFirstName() string {
	return u.firstname
}

func (u userInfo) GetUserName() string {
	return u.username
}

func (u userInfo) GetEmail() string {
	return u.email
}

func (u userInfo) GetId() string {
	return u.id
}
