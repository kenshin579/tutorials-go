package middleware

import (
	"net/http"
	"strings"

	"github.com/kenshin579/tutorials-go/web/sns-login/backend/service"
	"github.com/labstack/echo/v4"
)

// JWTAuth는 JWT 토큰을 검증하는 미들웨어
func JWTAuth(tokenService *service.TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization 헤더가 필요합니다")
			}

			// "Bearer <token>" 형식 파싱
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "잘못된 Authorization 형식")
			}

			claims, err := tokenService.ValidateToken(parts[1])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 토큰")
			}

			// Context에 사용자 ID 저장
			c.Set("user_id", claims.UserID)
			return next(c)
		}
	}
}
