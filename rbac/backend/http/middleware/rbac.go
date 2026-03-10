package middleware

import (
	"net/http"

	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"github.com/labstack/echo/v4"
)

// RequirePermission returns middleware that checks if the authenticated user
// has the specified permission.
func RequirePermission(permission string, permRepo domain.PermissionRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("user_id").(uint)

			permissions, err := permRepo.FindByUserID(userID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to check permissions")
			}

			for _, p := range permissions {
				if p.Key() == permission {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
}
