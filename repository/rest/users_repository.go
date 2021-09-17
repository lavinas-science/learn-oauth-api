package rest

import (
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lavinas-science/learn-oauth-api/domain/access_token"
	"github.com/lavinas-science/learn-oauth-api/domain/users"
	"github.com/lavinas-science/learn-oauth-api/utils/errors"
)

var (
	client = resty.New()
)

func init() {
	client.SetTimeout(1 * time.Second)
}

type RestUsersRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpires(access_token.AccessToken) *errors.RestErr
	LoginUser(string, string) (*users.User, *errors.RestErr)

}

type restUsersRepository struct {
}

func NewRepository() RestUsersRepository {
	return &restUsersRepository{}
}

func (r *restUsersRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	var restErr errors.RestErr
	var user users.User
	request := access_token.AccessTokenRequest{
		Username: email,
		Password: password,
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post("https://api.bookstore.com/users/login")
	if err != nil {
		return nil, errors.NewInternalServerError("Login api call error")
	}
	if resp.RawResponse.StatusCode > 299 {
		if err := json.Unmarshal(resp.Body(), &restErr); err != nil {
			return nil, errors.NewInternalServerError("invalid rest-client error unmarshall error")
		}
		return nil, &restErr
	}
	if err := json.Unmarshal(resp.Body(), &user); err != nil {
		return nil, errors.NewInternalServerError("Invalid rest-client error unmarshall client")
	}

	return &user, nil
}

func (r *restUsersRepository) GetById(string) (*access_token.AccessToken, *errors.RestErr) {
	return nil, errors.NewNotImplementedError("Not implemented")
}

func (r *restUsersRepository) Create(access_token.AccessToken) *errors.RestErr {
	return errors.NewNotImplementedError("Not implemented")
}

func (r *restUsersRepository) UpdateExpires(access_token.AccessToken) *errors.RestErr {
	return errors.NewNotImplementedError("Not implemented")
}