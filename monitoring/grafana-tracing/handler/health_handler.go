package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck 헬스체크 핸들러
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
