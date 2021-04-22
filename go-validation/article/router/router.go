package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Validator = NewValidator()
	return e
}
