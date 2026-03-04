package custom

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTConfig는 JWT 인증 미들웨어 설정
type JWTConfig struct {
	// Skipper는 미들웨어를 건너뛸 조건을 정의한다
	Skipper middleware.Skipper

	// SigningKey는 HMAC 서명 검증에 사용할 키
	SigningKey []byte

	// ContextKey는 검증된 클레임을 Context에 저장할 키 이름
	ContextKey string
}

// Claims는 JWT 토큰의 커스텀 클레임
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// DefaultJWTConfig는 기본 설정값
var DefaultJWTConfig = JWTConfig{
	Skipper:    middleware.DefaultSkipper,
	ContextKey: "user",
}

// JWTAuth는 JWT 토큰 검증 미들웨어를 반환한다
func JWTAuth(config JWTConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultJWTConfig.Skipper
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultJWTConfig.ContextKey
	}
	if len(config.SigningKey) == 0 {
		panic("jwt middleware: signing key is required")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// Authorization 헤더에서 Bearer 토큰 추출
			token, err := extractToken(c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			// 토큰 검증 및 클레임 파싱
			claims, err := validateToken(token, config.SigningKey)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 토큰")
			}

			// Context에 클레임 저장
			c.Set(config.ContextKey, claims)
			return next(c)
		}
	}
}

// extractToken는 Authorization 헤더에서 Bearer 토큰을 추출한다
func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization 헤더가 필요합니다")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("잘못된 Authorization 형식")
	}

	return parts[1], nil
}

// validateToken은 JWT 토큰을 검증하고 클레임을 반환한다
func validateToken(tokenString string, signingKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("잘못된 서명 방식")
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("유효하지 않은 클레임")
	}

	return claims, nil
}
