package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PingHandler struct{}

func NewPingHandler(e *echo.Echo) PingHandler {
	handler := PingHandler{}

	e.GET("/ping", handler.Ping)

	return handler
}

// Ping godoc
// @ID PingHandler.Ping
// @Tags Ping
// @Summary Ping API
// @Description Ping API
// @Success 200 {string} string "pong 메시지를 리턴한다"
// @Router /ping [get]
func (p PingHandler) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "Pong")
}
