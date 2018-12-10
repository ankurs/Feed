package cassandra

import (
	"context"
	"strings"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/ankurs/Feed/Feed/service/store/db"
	"github.com/carousell/Orion/utils/errors"
	"github.com/carousell/Orion/utils/log"
	"github.com/carousell/Orion/utils/spanutils"
	"github.com/gocql/gocql"
	"github.com/pborman/uuid"
)

type cas struct {
	casSes      *gocql.Session
	consistency gocql.Consistency
}

func (c *cas) CassandraExec(ctx context.Context, name, query string, values []interface{}, dest []interface{}) error {
	return c.CassandraExecWithConsistency(ctx, name, query, values, dest, c.consistency)
}

func (c *cas) CassandraExecWithConsistency(ctx context.Context, name, query string, values []interface{}, dest []interface{}, cons gocql.Consistency) error {
	// zipkin span
	span, ctx := spanutils.NewDatastoreSpan(ctx, name, "Cassandra")
	defer span.Finish()
	span.SetQuery(query)
	span.SetTag("values", values)

	var casError error
	e := hystrix.Do(name, func() error {
		q := c.casSes.Query(query, values...).Consistency(cons)
		if len(dest) == 0 {
			casError = q.Exec()
		} else {
			casError = q.Scan(dest...)
		}
		// don't count gocql.ErrNotFound in hystrix
		if casError == gocql.ErrNotFound {
			return nil
		}
		return casError
	}, nil)
	if e != nil {
		span.SetError(e.Error())
		return errors.Wrap(e, name)
	}
	if casError == gocql.ErrNotFound {
		return db.ErrNotFound
	}
	return casError
}

func (c *cas) AddFollowing(ctx context.Context, userId, followingId string) error {
	name := "AddFollowing"
	query := "INSERT INTO follow.following (user, following) VALUES (?,?)"

	return errors.Wrap(
		c.CassandraExec(
			ctx, name, query,
			db.BuildInterface(userId, followingId),
			db.BuildInterface(),
		),
		name)
}

func (c *cas) AddFollower(ctx context.Context, userId, followerId string) error {
	name := "AddFollower"
	query := "INSERT INTO follow.follower (user, follower) VALUES (?,?)"

	return errors.Wrap(
		c.CassandraExec(
			ctx, name, query,
			db.BuildInterface(userId, followerId),
			db.BuildInterface(),
		),
		name)
}

func (c *cas) RemoveFollowing(ctx context.Context, userId, followingId string) error {
	name := "RemoveFollowing"
	query := "DELETE FROM follow.following WHERE user= ? AND following = ?"

	return errors.Wrap(
		c.CassandraExec(
			ctx, name, query,
			db.BuildInterface(userId, followingId),
			db.BuildInterface(),
		),
		name)
}

func (c *cas) RemoveFollower(ctx context.Context, userId, followerId string) error {
	name := "RemoveFollower"
	query := "DELETE FROM follow.follower WHERE user = ? AND follower = ?"

	return errors.Wrap(
		c.CassandraExec(
			ctx, name, query,
			db.BuildInterface(userId, followerId),
			db.BuildInterface(),
		),
		name)
}

func (c *cas) CheckUserName(ctx context.Context, username string) (string, error) {
	name := "CheckUsername"
	query := "SELECT email FROM user.users WHERE username = ?"
	email := ""
	err := c.CassandraExec(
		ctx, name, query,
		db.BuildInterface(username),
		db.BuildInterface(&email),
	)
	return email, errors.Wrap(err, name)
}

func (c *cas) CheckEmail(ctx context.Context, email string) (string, error) {
	name := "CheckEmail"
	query := "SELECT email FROM user.users WHERE email = ?"
	mail := ""
	err := c.CassandraExec(
		ctx, name, query,
		db.BuildInterface(email),
		db.BuildInterface(&mail),
	)
	return mail, errors.Wrap(err, name)
}

func (c *cas) CheckLogin(ctx context.Context, username, password string, hash func(context.Context, string, string) string) (db.UserInfo, error) {
	username = strings.ToLower(username)
	name := "CheckLogin"
	query := "SELECT id,password,salt,email,firstname,lastname FROM user.users WHERE username = ?"
	id := ""
	pwd := ""
	salt := ""

	email := ""
	firstname := ""
	lastname := ""

	err := c.CassandraExec(
		ctx, name, query,
		db.BuildInterface(username),
		db.BuildInterface(&id, &pwd, &salt, &email, &firstname, &lastname),
	)
	if err == nil {
		log.Info(ctx, "salt", salt, "Password", password)
		if hash(ctx, password, salt) == pwd {
			return userInfo{
				email:     email,
				firstname: firstname,
				lastname:  lastname,
				username:  username,
				id:        id,
			}, nil
		} else {
			return nil, db.ErrNotFound
		}
	}
	return nil, errors.Wrap(err, name)
}
func (c *cas) CreateUser(ctx context.Context, req db.UserInfo, password string, hash func(context.Context, string, string) string) (string, error) {
	name := "CreateUser"

	id := uuid.New()
	salt := uuid.New() // TODO replace with crypto secure salt generation
	query := "INSERT INTO user.users (id, email, firstname, lastname, username, password, salt) VALUES (?,?,?,?,?,?,?)"

	password = hash(ctx, password, salt)
	err := c.CassandraExec(
		ctx, name, query,
		db.BuildInterface(id, req.GetEmail(), req.GetFirstName(), req.GetLastName(), strings.ToLower(req.GetUserName()), password, salt),
		db.BuildInterface(),
	)
	if err != nil {
		return "", errors.Wrap(err, name)
	}
	return id, nil
}

func (c *cas) GetUser(ctx context.Context, userID string) (db.UserInfo, error) {
	name := "GetUser"
	query := "SELECT username,email,firstname,lastname FROM user.users WHERE id= ?"

	username := ""
	email := ""
	firstname := ""
	lastname := ""

	err := c.CassandraExec(
		ctx, name, query,
		db.BuildInterface(userID),
		db.BuildInterface(&username, &email, &firstname, &lastname),
	)
	if err == nil {
		return userInfo{
			email:     email,
			firstname: firstname,
			lastname:  lastname,
			username:  username,
			id:        userID,
		}, nil
	}
	return nil, errors.Wrap(err, name)
}

func (c *cas) Close() {
	if c.casSes != nil {
		c.casSes.Close()
	}
}

func New(config Config) (Cassandra, error) {
	cluster := gocql.NewCluster(config.Hosts...)
	ses, err := cluster.CreateSession()
	if err != nil {
		log.Error(context.Background(), err)
		return nil, errors.Wrap(err, "New Cassandra connection")
	}
	return &cas{
		casSes:      ses,
		consistency: gocql.LocalQuorum,
	}, nil
}
