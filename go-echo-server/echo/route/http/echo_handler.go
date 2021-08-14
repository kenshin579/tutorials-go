package http

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo"
)

type echoHandler struct {
}

func NewEchoHandler(e *echo.Echo) *echoHandler {
	handler := &echoHandler{}

	e.POST("/", handler.EchoHandler)

	return handler
}

func (w *echoHandler) EchoHandler(ctx echo.Context) error {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error(err)
	}
	log.Infof("request:%v", string(body))

	return ctx.NoContent(http.StatusOK)
}
