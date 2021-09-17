package rest

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/lavinas-science/learn-oauth-api/domain/access_token"
	"github.com/lavinas-science/learn-oauth-api/domain/users"
	"github.com/lavinas-science/learn-oauth-api/utils/errors"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()
	os.Exit(m.Run())
}

func TestLoginTimeoutFromApi(t *testing.T) {
	// Mock response
	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
		func(req *http.Request) (*http.Response, error) {
			var lg access_token.AccessTokenRequest
			if err := json.NewDecoder(req.Body).Decode(&lg); err != nil {
				return httpmock.NewJsonResponse(http.StatusBadRequest, "Bad request")
			}
			us := users.User{Id: 100, FirstName: "First Name", LastName: lg.Password, Email: lg.Username}
			resp, err := httpmock.NewJsonResponse(http.StatusOK, us)
			if err != nil {
				return httpmock.NewJsonResponse(http.StatusInternalServerError, "Internal error")
			}

			time.Sleep(2 * time.Second)

			return resp, nil
		},
	)
	rep := NewRepository()
	user, err := rep.LoginUser("login@user.com", "passwd")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.Contains(t, err.Message, "Login api call error")
}

func TestLoginInvalidErrorInterface(t *testing.T) {
	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
		httpmock.NewStringResponder(http.StatusInternalServerError, ``))
	rep := NewRepository()
	user, err := rep.LoginUser("login@user.com", "passwd")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid rest-client error unmarshall error", err.Message)
}

func TestLoginInvalidErrorInterface2(t *testing.T) {
	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
		httpmock.NewStringResponder(http.StatusInternalServerError, `"xxx": 123`))
	rep := NewRepository()
	user, err := rep.LoginUser("login@user.com", "passwd")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid rest-client error unmarshall error", err.Message)
}

func TestLoginInvalidUserInterface(t *testing.T) {
	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
		httpmock.NewStringResponder(http.StatusOK, ``))
	rep := NewRepository()
	user, err := rep.LoginUser("login@user.com", "passwd")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "Invalid rest-client error unmarshall client", err.Message)
}

func TestLoginInvalidLoginCredentials(t *testing.T) {
	// Mock response
	errRest := errors.RestErr{Status: http.StatusNotFound, Error: "not_found", Message: "No record found"}
	resp, _ := httpmock.NewJsonResponder(http.StatusNotFound, errRest)
	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login", resp)
	rep := NewRepository()
	user, err := rep.LoginUser("login@user.com", "passwd")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, errRest.Status, err.Status)
	assert.EqualValues(t, errRest.Message, err.Message)
	assert.EqualValues(t, errRest.Error, err.Error)

}

func TestLoginOk(t *testing.T) {
	// Mock response
	httpmock.RegisterResponder("POST", "https://api.bookstore.com/users/login",
		func(req *http.Request) (*http.Response, error) {
			var lg access_token.AccessTokenRequest
			if err := json.NewDecoder(req.Body).Decode(&lg); err != nil {
				return httpmock.NewJsonResponse(http.StatusBadRequest, "Bad request")
			}
			us := users.User{Id: 100, FirstName: "First Name", LastName: lg.Password, Email: lg.Username}
			resp, err := httpmock.NewJsonResponse(http.StatusOK, us)
			if err != nil {
				return httpmock.NewJsonResponse(http.StatusInternalServerError, "Internal error")
			}
			return resp, nil
		},
	)
	// call func
	rep := NewRepository()
	user, err := rep.LoginUser("login@user.com", "passwd")
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, 100, user.Id)
	assert.EqualValues(t, "First Name", user.FirstName)
	assert.EqualValues(t, "passwd", user.LastName)
	assert.EqualValues(t, "login@user.com", user.Email)
}
