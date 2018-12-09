package store

import (
	"context"
	"strings"

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
	ErrInvalidLogin = errors.New("could not login to the account")
)

func (s *str) Register(ctx context.Context, req RegisterRequest) (LoginResponse, error) {
	name := "StorageRegister"
	// zipkin span
	span, ctx := spanutils.NewInternalSpan(ctx, name)
	defer span.Finish()

	username := strings.ToLower(req.GetUserName())
	_, err := s.checkUserName(ctx, username)
	if cause(err) != gocql.ErrNotFound {
		if err == nil {
			return nil, ErrAlreadyTaken
		}
		return nil, errors.Wrap(err, name)
	}

	_, err = s.checkEmail(ctx, req.GetEmail())
	if cause(err) != gocql.ErrNotFound {
		if err == nil {
			return nil, ErrAlreadyTaken
		}
		return nil, errors.Wrap(err, name)
	}

	id, err := s.createUser(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, name)
	}
	resp := loginResponse{
		userInfo: userInfo{
			firstname: req.GetFirstName(),
			lastname:  req.GetLastName(),
			email:     req.GetEmail(),
			id:        id,
			username:  strings.ToLower(req.GetUserName()),
		},
	}
	return resp, nil
}

func (s *str) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {

	user, err := s.checkLogin(ctx, req.GetUserName(), req.GetPassword())
	if err == nil && user != nil {
		// for now just use id as login token
		// TODO move to JWT token
		return loginResponse{token: user.GetId(), userInfo: user}, nil
	}
	return nil, errors.Wrap(err, "Login")
}

func (s *str) checkLogin(ctx context.Context, username, password string) (UserInfo, error) {
	username = strings.ToLower(username)
	name := "CheckLogin"
	query := "SELECT id,password,salt,email,firstname,lastname FROM user.users WHERE username = ?"
	id := ""
	pwd := ""
	salt := ""

	email := ""
	firstname := ""
	lastname := ""

	err := s.cassandraExec(
		ctx, name, query,
		buildInterface(username),
		buildInterface(&id, &pwd, &salt, &email, &firstname, &lastname),
	)
	if err == nil {
		log.Info(ctx, "salt", salt, "Password", password)
		if getPasswordHash(ctx, password, salt) == pwd {
			return userInfo{
				email:     email,
				firstname: firstname,
				lastname:  lastname,
				username:  username,
				id:        id,
			}, nil
		} else {
			return nil, ErrInvalidLogin
		}
	}
	if err == gocql.ErrNotFound {
		return nil, ErrInvalidLogin
	}
	return nil, errors.Wrap(err, name)
}

func (s *str) createUser(ctx context.Context, req RegisterRequest) (string, error) {
	name := "CreateUser"

	id := uuid.New()
	salt := uuid.New() // TODO replace with crypto secure salt generation
	query := "INSERT INTO user.users (id, email, firstname, lastname, username, password, salt) VALUES (?,?,?,?,?,?,?)"

	password := getPasswordHash(ctx, req.GetPassword(), salt)
	err := s.cassandraExec(
		ctx, name, query,
		buildInterface(id, req.GetEmail(), req.GetFirstName(), req.GetLastName(), strings.ToLower(req.GetUserName()), password, salt),
		buildInterface(),
	)
	if err != nil {
		return id, errors.Wrap(err, name)
	}
	return id, nil
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
