package http

import (
	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"github.com/kenshin579/tutorials-go/rbac/backend/http/handler"
	"github.com/kenshin579/tutorials-go/rbac/backend/http/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Handlers holds all HTTP handlers.
type Handlers struct {
	Auth    *handler.AuthHandler
	User    *handler.UserHandler
	Rbac    *handler.RbacHandler
	Product *handler.ProductHandler
	Order   *handler.OrderHandler
}

// SetupRoutes configures all API routes with middleware chains.
func SetupRoutes(e *echo.Echo, h *Handlers, jwtSecret string, permRepo domain.PermissionRepository, db *gorm.DB) {
	api := e.Group("/api")

	// Auth routes (no auth middleware)
	auth := api.Group("/auth")
	auth.POST("/register", h.Auth.Register)
	auth.POST("/login", h.Auth.Login)
	auth.POST("/refresh", h.Auth.Refresh)
	auth.POST("/logout", h.Auth.Logout, middleware.JWTAuth(jwtSecret))

	// Secured routes (require JWT)
	jwtMw := middleware.JWTAuth(jwtSecret)
	secured := api.Group("", jwtMw)

	// Helper to create RBAC middleware
	rbac := func(permission string) echo.MiddlewareFunc {
		return middleware.RequirePermission(permission, permRepo)
	}

	// Helper to create Owner middleware
	owner := func(config middleware.OwnerConfig) echo.MiddlewareFunc {
		return middleware.RequireOwner(config, db)
	}

	// Products
	products := secured.Group("/products")
	products.GET("", h.Product.List, rbac("products:read"))
	products.GET("/:id", h.Product.Get, rbac("products:read"))
	products.POST("", h.Product.Create, rbac("products:create"))
	products.PUT("/:id", h.Product.Update, rbac("products:update"),
		owner(middleware.OwnerConfig{ResourceTable: "products", OwnerField: "created_by", BypassRoles: []string{"admin"}}))
	products.DELETE("/:id", h.Product.Delete, rbac("products:delete"))
	products.PATCH("/:id/status", h.Product.UpdateStatus, rbac("products:status:update"))

	// Orders
	orders := secured.Group("/orders")
	orders.GET("", h.Order.List, rbac("orders:read"))
	orders.GET("/:id", h.Order.Get, rbac("orders:read"),
		owner(middleware.OwnerConfig{ResourceTable: "orders", OwnerField: "ordered_by", BypassRoles: []string{"admin", "manager"}}))
	orders.POST("", h.Order.Create, rbac("orders:create"))
	orders.PATCH("/:id/status", h.Order.UpdateStatus, rbac("orders:status:update"))
	orders.PATCH("/:id/cancel", h.Order.Cancel, rbac("orders:cancel"),
		owner(middleware.OwnerConfig{ResourceTable: "orders", OwnerField: "ordered_by", BypassRoles: []string{"admin", "manager"}}))

	// Users
	users := secured.Group("/users")
	users.GET("", h.User.List, rbac("users:read"))
	users.GET("/:id", h.User.Get, rbac("users:read"))
	users.PUT("/:id", h.User.Update, rbac("users:update"))
	users.DELETE("/:id", h.User.Delete, rbac("users:delete"))
	users.POST("/:id/roles", h.User.AssignRole, rbac("users:update"))
	users.DELETE("/:id/roles/:roleId", h.User.RemoveRole, rbac("users:update"))

	// Roles
	roles := secured.Group("/roles")
	roles.GET("", h.Rbac.ListRoles, rbac("roles:read"))
	roles.POST("", h.Rbac.CreateRole, rbac("roles:create"))
	roles.PUT("/:id", h.Rbac.UpdateRole, rbac("roles:update"))
	roles.DELETE("/:id", h.Rbac.DeleteRole, rbac("roles:delete"))
	roles.POST("/:id/permissions", h.Rbac.AssignPermission, rbac("roles:update"))
	roles.DELETE("/:id/permissions/:permId", h.Rbac.RemovePermission, rbac("roles:update"))

	// Permissions
	secured.GET("/permissions", h.Rbac.ListPermissions, rbac("roles:read"))
}
