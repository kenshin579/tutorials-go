package main

import (
	"keycloak-tutorial-backend/handlers"
	"keycloak-tutorial-backend/middleware"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	// CORS middleware for local development
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Protected routes
	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware())
	api.GET("/user", handlers.GetUserInfo)

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}
