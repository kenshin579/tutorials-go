package router

import (
	"github.com/kenshin579/tutorials-go/go-redis/blackboard/common/errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))
	e.Validator = NewValidator()
	e.HTTPErrorHandler = errors.HTTPErrorHandler
	return e
}
