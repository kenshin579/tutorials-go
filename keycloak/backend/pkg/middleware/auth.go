package middleware

import (
	"net/http"
	"strings"

	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/domain"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(userUseCase domain.UserUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := extractToken(c)
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
			}

			ctx := c.Request().Context()
			valid, err := userUseCase.ValidateToken(ctx, token)
			if err != nil || !valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			// 토큰을 컨텍스트에 저장
			c.Set("token", token)

			return next(c)
		}
	}
}

func extractToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return authHeader
}
