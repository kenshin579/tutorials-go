// Package server wires the Echo HTTP server with middleware and route registration.
// It is intentionally domain-agnostic — domain logic is in the todo package.
package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/kenshin579/tutorials-go/ai/superpowers/todo/backend/todo"
)

// New constructs a configured *echo.Echo bound to the given store.
// Middleware order: Recover → Logger → CORS. CORS allows the Vite dev/preview origins.
func New(store *todo.Store) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:4173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
	}))

	h := todo.NewHandler(store)

	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	e.GET("/api/todos", h.List)
	e.POST("/api/todos", h.Create)
	e.PATCH("/api/todos/:id", h.Update)
	e.DELETE("/api/todos/:id", h.Delete)

	return e
}
