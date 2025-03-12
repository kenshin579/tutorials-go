package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Validator = NewValidator()
	return e
}
