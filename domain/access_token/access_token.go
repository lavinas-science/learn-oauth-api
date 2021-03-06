package access_token

import (
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"github.com/lavinas-science/learn-utils-go/rest_errors"
)

const (
	expirationTime             = 24
	grandTypePassword          = "password"
	grandTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType    string `json:"grand_type"`
	Scope        string `json:"scope"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() rest_errors.RestErr {

	switch at.GrantType {
	case grandTypePassword:
		break
	case grandTypeClientCredentials:
		break
	default:
		return rest_errors.NewBadRequestError(" ")
	}
	return nil
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (at *AccessToken) Validate() rest_errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return rest_errors.NewBadRequestError("invalid access token id")
	}
	if at.ClientId <= 0 {
		return rest_errors.NewBadRequestError("invalid client id")
	}
	if at.UserId <= 0 {
		return rest_errors.NewBadRequestError("invalid user id")
	}
	if at.Expires <= 0 {
		return rest_errors.NewBadRequestError("invalid expiration time")
	}
	return nil

}

func GetNewAccessToken(user_id int64) AccessToken {
	b := make([]byte, tokenLength)
	rand.Read(b)
	token := hex.EncodeToString(b)
	return AccessToken{
		AccessToken: token,
		UserId:      user_id,
		ClientId:    user_id,
		Expires:     time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}
