package rest

import (
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lavinas-science/learn-oauth-api/domain/access_token"
	"github.com/lavinas-science/learn-oauth-api/domain/users"
	"github.com/lavinas-science/learn-oauth-api/utils/errors"
)


const (
	UserContentType = "application/json"
	UserBaseURI = "http://127.0.0.1:8080"
	UserURI = "/users/login"
	PingURI = "/ping"
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
	Ping() bool
}

type restUsersRepository struct {
}

type clientRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func NewRepository() RestUsersRepository {
	return &restUsersRepository{}
}

func (r *restUsersRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	var restErr errors.RestErr
	var user users.User
	request := clientRequest{Email: email, Password: password}
	resp, err := client.R().
		SetHeader("Content-Type",UserContentType).
		SetBody(request).
		Post(UserBaseURI + UserURI)
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

func (r *restUsersRepository) Ping() bool {
	url := UserBaseURI + PingURI
	resp, err := client.R().Execute("GET", url)
	if err != nil {
		return false
	}
	if resp.RawResponse.StatusCode != 200 {
		return false
	}
	return true
} 