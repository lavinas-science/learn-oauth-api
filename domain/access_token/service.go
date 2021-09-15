package access_token

import (
	"strings"
	
	"github.com/lavinas-science/learn-oauth-api/utils/errors"
)

type Repository interface {
	GetById(string) (*AccessToken, *errors.RestErr)
	Create(AccessToken) *errors.RestErr
	UpdateExpires(AccessToken) *errors.RestErr
}

type Service interface{
	GetById(string) (*AccessToken, *errors.RestErr)
	Create(AccessToken) *errors.RestErr
	UpdateExpires(AccessToken) *errors.RestErr
}

type service struct {	
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

func (s *service) GetById(accessTokenId string) (*AccessToken, *errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, errors.NewBadRequestError("invalid access token")
	}
	return s.repository.GetById(accessTokenId)
} 

func (s *service) Create(AccessToken) *errors.RestErr {
	return nil
}

func (s *service) UpdateExpires(AccessToken) *errors.RestErr {
	return nil
}
