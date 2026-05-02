package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueAndParse(t *testing.T) {
	secret := "test-secret"
	token, err := Issue(42, secret, time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := Parse(token, secret)
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
}

func TestParse_Expired(t *testing.T) {
	secret := "test-secret"
	token, _ := Issue(1, secret, -time.Hour)
	_, err := Parse(token, secret)
	assert.Error(t, err)
}

func TestParse_WrongSecret(t *testing.T) {
	token, _ := Issue(1, "secret-a", time.Hour)
	_, err := Parse(token, "secret-b")
	assert.Error(t, err)
}
