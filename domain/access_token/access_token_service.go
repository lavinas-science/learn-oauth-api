package access_token

import (
	"strings"

	"github.com/lavinas-science/learn-oauth-api/domain/users"
	"github.com/lavinas-science/learn-utils-go/rest_errors"
)

const (
	tokenLength = 32
)

type Repository interface {
	GetById(string) (*AccessToken, rest_errors.RestErr)
	Create(AccessToken) rest_errors.RestErr
	UpdateExpires(AccessToken) rest_errors.RestErr
	LoginUser(string, string) (*users.User, rest_errors.RestErr)
	Ping() bool
}

type Service interface {
	GetById(string) (*AccessToken, rest_errors.RestErr)
	Create(AccessTokenRequest) (*AccessToken, rest_errors.RestErr)
	UpdateExpires(AccessToken) rest_errors.RestErr
}

type service struct {
	db_repository   Repository
	user_repository Repository
}

func NewService(db_repo Repository, user_repo Repository) Service {
	return &service{
		db_repository:   db_repo,
		user_repository: user_repo,
	}
}

func (s *service) GetById(accessTokenId string) (*AccessToken, rest_errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token")
	}
	return s.db_repository.GetById(accessTokenId)
}

func (s *service) Create(atr AccessTokenRequest) (*AccessToken, rest_errors.RestErr) {
	if err := atr.Validate(); err != nil {
		return nil, err
	}
	user, err := s.user_repository.LoginUser(atr.Username, atr.Password)
	if err != nil {
		return nil, err
	}
	at := GetNewAccessToken(user.Id)
	if err := s.db_repository.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpires(at AccessToken) rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return nil
}
