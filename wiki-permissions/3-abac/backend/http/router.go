package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/http/handler"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/http/middleware"
)

type Deps struct {
	JWTSecret  string
	Auth       *handler.AuthHandler
	Page       *handler.PageHandler
	Department *handler.DepartmentHandler
}

// NewRouter는 Echo 인스턴스를 구성하고 모든 라우트를 등록한다.
// /auth/login 외 모든 /api 라우트는 JWT 미들웨어로 보호되며,
// ABAC 정책 평가는 미들웨어가 아니라 usecase 단에서 수행한다 (1·2편과 동일한 디자인).
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

	api.GET("/departments", d.Department.List)

	return e
}
