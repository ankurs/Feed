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

type feedInfo struct {
	id     string
	actor  string
	verb   string
	cverb  string
	object string
	target string
	ts     int64
}

func (f feedInfo) GetId() string {
	return f.id
}

func (f feedInfo) GetActor() string {
	return f.actor
}

func (f feedInfo) GetVerb() string {
	return f.verb
}

func (f feedInfo) GetCVerb() string {
	return f.cverb
}

func (f feedInfo) GetObject() string {
	return f.object
}

func (f feedInfo) GetTarget() string {
	return f.target
}

func (f feedInfo) GetTs() int64 {
	return f.ts
}
