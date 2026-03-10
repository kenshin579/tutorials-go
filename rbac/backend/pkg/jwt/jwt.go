package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims for this application.
type Claims struct {
	UserID uint     `json:"user_id"`
	Roles  []string `json:"roles"`
	jwtv5.RegisteredClaims
}

// TokenPair holds access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateTokenPair creates an access token (15 min) and refresh token (7 days).
func GenerateTokenPair(userID uint, roles []string, secret string) (*TokenPair, error) {
	// Access Token: 15 minutes
	accessClaims := &Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	}
	accessToken, err := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, accessClaims).SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	// Refresh Token: 7 days
	refreshClaims := &Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ParseToken validates and parses a JWT token string.
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(tokenString, &Claims{}, func(token *jwtv5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
