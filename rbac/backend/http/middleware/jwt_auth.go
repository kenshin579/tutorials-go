package middleware

import (
	"net/http"
	"strings"

	"github.com/kenshin579/tutorials-go/rbac/backend/pkg/jwt"
	"github.com/labstack/echo/v4"
)

// JWTAuth returns middleware that validates JWT tokens from the Authorization header.
func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := extractBearerToken(c.Request().Header.Get("Authorization"))
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}

			claims, err := jwt.ParseToken(token, jwtSecret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			c.Set("user_id", claims.UserID)
			c.Set("roles", claims.Roles)
			return next(c)
		}
	}
}

func extractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
