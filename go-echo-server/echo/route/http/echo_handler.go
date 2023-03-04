package http

import (
	"io"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
)

type echoHandler struct {
}

func NewEchoHandler(e *echo.Echo) *echoHandler {
	handler := &echoHandler{}

	e.POST("/echo", handler.EchoHandler)

	return handler
}

// EchoHandler godoc
// @ID EchoHandler
// @Tags Echo
// @Summary Echo API
// @Description Echo API
// @Success 200 {string} Body "Request 값을 그대로 반환한다"
// @Router /echo [get]
func (handler echoHandler) EchoHandler(ctx echo.Context) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error(err)
	}
	log.Infof("request:%v", string(body))

	return ctx.JSON(http.StatusOK, body)
}
