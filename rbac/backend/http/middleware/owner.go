package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// OwnerConfig defines configuration for the owner middleware.
type OwnerConfig struct {
	ResourceTable string   // DB table name (e.g. "products")
	OwnerField    string   // Owner column name (e.g. "created_by")
	BypassRoles   []string // Roles that skip the owner check (e.g. ["admin"])
}

// RequireOwner returns middleware that checks if the authenticated user
// is the owner of the requested resource, unless they have a bypass role.
func RequireOwner(config OwnerConfig, db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("user_id").(uint)
			roles := c.Get("roles").([]string)

			// Check bypass roles
			for _, role := range roles {
				for _, bypassRole := range config.BypassRoles {
					if role == bypassRole {
						return next(c)
					}
				}
			}

			// Extract resource ID from path parameter
			resourceID := c.Param("id")

			// Query DB for the owner field
			var ownerID uint
			err := db.Table(config.ResourceTable).
				Where("id = ?", resourceID).
				Pluck(config.OwnerField, &ownerID).Error
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, "resource not found")
			}

			if ownerID != userID {
				return echo.NewHTTPError(http.StatusForbidden, "not the owner of this resource")
			}

			return next(c)
		}
	}
}
