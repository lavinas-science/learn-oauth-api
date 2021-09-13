package access_token

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)


func TestAccessTokenConstants (t *testing.T) {
	assert.EqualValues(t, 24, expirationTime, "expiration time should by 24 hours")
}

func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken()
	// assert.False(t, at.IsExpired(), "new accesstoken should not be expired")
	assert.EqualValues(t, at.AccessToken, "", "new access token should not have defined id")
	assert.EqualValues(t, at.UserId, 0, "new access token should not have defined user_id")
} 

func TestAccessTokenIsExpired(t *testing.T) {
	at := AccessToken{}
	// assert.True(t, at.IsExpired(), "accesstoken should be expired")
	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "accesstoken expriring at more 3 housrs should not be expired")
	at.Expires = time.Now().UTC().Add(-1 * time.Hour).Unix()
	// assert.True(t, at.IsExpired(), "accesstoken expired at 1 hour ago should be expired")
}