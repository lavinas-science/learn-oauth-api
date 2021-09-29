package rest

import (
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lavinas-science/learn-oauth-api/domain/access_token"
	"github.com/lavinas-science/learn-oauth-api/domain/users"
	"github.com/lavinas-science/learn-utils-go/rest_errors"
)

const (
	UserContentType = "application/json"
	UserBaseURI     = "http://127.0.0.1:8080"
	UserURI         = "/users/login"
	PingURI         = "/ping"
	timeoutSeconds  = 1
)

var (
	client = resty.New()
)

func init() {
	client.SetTimeout(timeoutSeconds * time.Second)
}

type RestUsersRepository interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessToken) rest_errors.RestErr
	UpdateExpires(access_token.AccessToken) rest_errors.RestErr
	LoginUser(string, string) (*users.User, rest_errors.RestErr)
	Ping() bool
}

type restUsersRepository struct {
}

type clientRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRepository() RestUsersRepository {
	return &restUsersRepository{}
}

func (r *restUsersRepository) LoginUser(email string, password string) (*users.User, rest_errors.RestErr) {
	var restErr rest_errors.RestErr
	var user users.User
	request := clientRequest{Email: email, Password: password}
	resp, err := client.R().
		SetHeader("Content-Type", UserContentType).
		SetBody(request).
		Post(UserBaseURI + UserURI)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("Login api call error")
	}
	if resp.RawResponse.StatusCode > 299 {
		if restErr, err = rest_errors.NewRestErrorFromBytes(resp.Body()); err != nil {
			return nil, rest_errors.NewInternalServerError("invalid rest-client error unmarshall error")
		}
		return nil, restErr
	}
	if err := json.Unmarshal(resp.Body(), &user); err != nil {
		return nil, rest_errors.NewInternalServerError("invalid rest-client error unmarshall client")
	}
	return &user, nil
}

func (r *restUsersRepository) GetById(string) (*access_token.AccessToken, rest_errors.RestErr) {
	return nil, rest_errors.NewNotImplementedError("Not implemented")
}

func (r *restUsersRepository) Create(access_token.AccessToken) rest_errors.RestErr {
	return rest_errors.NewNotImplementedError("Not implemented")
}

func (r *restUsersRepository) UpdateExpires(access_token.AccessToken) rest_errors.RestErr {
	return rest_errors.NewNotImplementedError("Not implemented")
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
