package jwt

import (
	"testing"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key"

func TestGenerateTokenPair(t *testing.T) {
	pair, err := GenerateTokenPair(1, []string{"admin"}, testSecret)
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.NotEqual(t, pair.AccessToken, pair.RefreshToken)
}

func TestParseToken(t *testing.T) {
	pair, err := GenerateTokenPair(42, []string{"admin", "manager"}, testSecret)
	require.NoError(t, err)

	claims, err := ParseToken(pair.AccessToken, testSecret)
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, []string{"admin", "manager"}, claims.Roles)
}

func TestParseToken_InvalidSecret(t *testing.T) {
	pair, err := GenerateTokenPair(1, []string{"user"}, testSecret)
	require.NoError(t, err)

	_, err = ParseToken(pair.AccessToken, "wrong-secret")
	assert.Error(t, err)
}

func TestParseToken_Expired(t *testing.T) {
	// Create an expired token manually
	claims := &Claims{
		UserID: 1,
		Roles:  []string{"user"},
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token, err := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	require.NoError(t, err)

	_, err = ParseToken(token, testSecret)
	assert.Error(t, err)
}

func TestParseRefreshToken(t *testing.T) {
	pair, err := GenerateTokenPair(5, []string{"user"}, testSecret)
	require.NoError(t, err)

	claims, err := ParseToken(pair.RefreshToken, testSecret)
	require.NoError(t, err)
	assert.Equal(t, uint(5), claims.UserID)
}
