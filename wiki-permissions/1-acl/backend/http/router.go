package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/http/handler"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/http/middleware"
)

// Deps는 NewRouter가 라우트에 와이어할 의존성 묶음이다.
type Deps struct {
	JWTSecret string
	Auth      *handler.AuthHandler
	Page      *handler.PageHandler
	ACL       *handler.ACLHandler
}

// NewRouter는 Echo 인스턴스를 구성하고 모든 라우트를 등록한다.
// /auth/login 외 모든 /api 라우트는 JWT 미들웨어로 보호된다.
func NewRouter(d Deps) *echo.Echo {
	e := echo.New()
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.POST("/auth/login", d.Auth.Login)

	api := e.Group("/api", middleware.JWTAuth(d.JWTSecret))
	api.GET("/pages", d.Page.List)
	api.GET("/pages/:id", d.Page.Get)
	api.PUT("/pages/:id", d.Page.Update)
	api.GET("/pages/:id/acl", d.ACL.List)
	api.POST("/pages/:id/acl", d.ACL.Grant)
	api.DELETE("/pages/:id/acl/:userId", d.ACL.Revoke)

	return e
}
