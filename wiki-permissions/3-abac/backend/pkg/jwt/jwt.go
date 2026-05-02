package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// Claims는 wiki-permissions 인증에 사용하는 JWT 클레임이다.
type Claims struct {
	UserID uint `json:"user_id"`
	jwtv5.RegisteredClaims
}

// Issue는 user_id를 담은 access token을 secret으로 서명하여 발급한다.
func Issue(userID uint, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwtv5.NewNumericDate(now),
		},
	}
	tok := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secret))
}

// Parse는 token 문자열을 secret으로 검증하고 Claims 포인터를 반환한다.
func Parse(tokenStr, secret string) (*Claims, error) {
	tok, err := jwtv5.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtv5.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	c, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return c, nil
}
