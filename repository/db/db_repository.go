package db

import (
	"github.com/gocql/gocql"
	"github.com/lavinas-science/learn-oauth-api/clients/cassandra"
	"github.com/lavinas-science/learn-oauth-api/domain/access_token"
	"github.com/lavinas-science/learn-oauth-api/domain/users"
	"github.com/lavinas-science/learn-utils-go/rest_errors"
)

const (
	getAccessToken    = "select access_token, user_id, client_id, expires from access_tokens where access_token = ?;"
	createAccessToken = "insert into access_tokens(access_token, user_id, client_id, expires) values (?, ?, ?, ?);"
	updateAccessToken = "update access_tokens set expires = ? where access_token = ?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create(access_token.AccessToken) *rest_errors.RestErr
	UpdateExpires(access_token.AccessToken) *rest_errors.RestErr
	LoginUser(string, string) (*users.User, *rest_errors.RestErr)
	Ping() bool
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *rest_errors.RestErr) {
	session := cassandra.GetSession()
	var at access_token.AccessToken
	if err := session.Query(getAccessToken, id).Scan(
		&at.AccessToken, &at.UserId, &at.ClientId, &at.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError("no access token found with given id")
		}
		return nil, rest_errors.NewInternalServerError(err.Error())

	}
	return &at, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *rest_errors.RestErr {
	session := cassandra.GetSession()
	if err := session.Query(
		createAccessToken, at.AccessToken, at.UserId,
		at.ClientId, at.Expires).Exec(); err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *dbRepository) UpdateExpires(at access_token.AccessToken) *rest_errors.RestErr {
	session := cassandra.GetSession()
	if err := session.Query(
		updateAccessToken, at.Expires, at.AccessToken).Exec(); err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	return nil
}


func (r *dbRepository) LoginUser(string, string) (*users.User, *rest_errors.RestErr) {
	return nil, rest_errors.NewNotImplementedError("Not implemented")
}

func (r *dbRepository) Ping() bool {
	return false
}
