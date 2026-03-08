package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/api/ping", Ping)
	e.GET("/api/resource", GetResource)
}

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "pong",
	})
}

func GetResource(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":   1,
		"name": "sample resource",
	})
}
