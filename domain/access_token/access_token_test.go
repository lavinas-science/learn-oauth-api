package access_token

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccessTokenConstants(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime, "expiration time should by 24 hours")
}

func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken(100)
	assert.EqualValues(t, at.UserId, 100)
	assert.EqualValues(t, at.ClientId, 100)
	assert.Greater(t, at.Expires, time.Now().UTC().Add(expirationTime * time.Hour - 10 * time.Minute).Unix())
	assert.Less(t, at.Expires, time.Now().UTC().Add(expirationTime * time.Hour + 10 * time.Minute).Unix())
}

func TestAccessTokenIsExpired(t *testing.T) {
	at := AccessToken{}
	// assert.True(t, at.IsExpired(), "accesstoken should be expired")
	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "accesstoken expriring at more 3 housrs should not be expired")
	at.Expires = time.Now().UTC().Add(-1 * time.Hour).Unix()
	// assert.True(t, at.IsExpired(), "accesstoken expired at 1 hour ago should be expired")
}
