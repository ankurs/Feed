package store

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/carousell/Orion/utils/errors"
	"github.com/carousell/Orion/utils/log"
	"github.com/carousell/Orion/utils/spanutils"
	"github.com/gocql/gocql"
	"github.com/pborman/uuid"
)

type str struct {
	casSes      *gocql.Session
	consistency gocql.Consistency
}

var (
	ErrAlreadyTaken = errors.New("error already taken")
)

type registerResponse struct {
	id string
}

func (r registerResponse) GetId() string {
	return r.id
}

func (s *str) Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	name := "StorageRegister"
	// zipkin span
	span, ctx := spanutils.NewInternalSpan(ctx, name)
	defer span.Finish()

	_, err := s.checkUserName(ctx, req.GetUserName())
	if err != gocql.ErrNotFound {
		if err == nil {
			return nil, ErrAlreadyTaken
		}
		return nil, errors.Wrap(err, name)
	}

	_, err = s.checkEmail(ctx, req.GetEmail())
	if err != gocql.ErrNotFound {
		if err == nil {
			return nil, ErrAlreadyTaken
		}
		return nil, errors.Wrap(err, name)
	}

	id, err := s.createUser(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, name)
	}
	resp := registerResponse{id: id}
	return resp, nil
}

func (s *str) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {

	return nil, nil
}

func buildInterface(vals ...interface{}) []interface{} {
	return vals
}

func (s *str) createUser(ctx context.Context, req RegisterRequest) (string, error) {
	name := "CreateUser"

	uuid := uuid.New()
	query := "INSERT INTO user.users (id, email, firstname, lastname, password, username) VALUES (?,?,?,?,?,?)"

	err := s.cassandraExec(
		ctx, name, query,
		buildInterface(uuid, req.GetEmail(), req.GetFirstName(), req.GetLastName(), req.GetPassword(), req.GetUserName()),
		buildInterface(),
	)
	if err != nil {
		return uuid, errors.Wrap(err, name)
	}
	return uuid, nil
}

func (s *str) checkUserName(ctx context.Context, username string) (string, error) {
	name := "CheckUsername"
	query := "SELECT email FROM user.users WHERE username = ?"
	email := ""
	err := s.cassandraExec(
		ctx, name, query,
		buildInterface(username),
		buildInterface(&email),
	)
	return email, errors.Wrap(err, name)
}

func (s *str) checkEmail(ctx context.Context, email string) (string, error) {
	name := "CheckEmail"
	query := "SELECT email FROM user.users WHERE email = ?"
	mail := ""
	err := s.cassandraExec(
		ctx, name, query,
		buildInterface(email),
		buildInterface(&mail),
	)
	return mail, errors.Wrap(err, name)
}

func (s *str) cassandraExec(ctx context.Context, name, query string, values []interface{}, dest []interface{}) error {
	// zipkin span
	span, ctx := spanutils.NewDatastoreSpan(ctx, name, "Cassandra")
	defer span.Finish()
	span.SetQuery(query)
	span.SetTag("values", values)

	var casError error
	e := hystrix.Do(name, func() error {
		q := s.casSes.Query(query, values...).Consistency(s.consistency)
		if len(dest) == 0 {
			casError = q.Exec()
		} else {
			casError = q.Scan(dest...)
		}
		if casError == gocql.ErrNotFound {
			return nil
		}
		return casError
	}, nil)
	if e != nil {
		span.SetError(e.Error())
		return errors.Wrap(e, name)
	}
	return casError
}

func (s *str) Close() {
	if s != nil && s.casSes != nil {
		s.casSes.Close()
	}
}

func NewStore(config Config) (Storage, error) {
	cluster := gocql.NewCluster(config.CassandraHosts...)
	ses, err := cluster.CreateSession()
	if err != nil {
		log.Error(context.Background(), err)
		return nil, errors.Wrap(err, "NewStore")
	}
	return &str{
		casSes:      ses,
		consistency: gocql.LocalQuorum,
	}, nil
}
