package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", helloHandler)
	e.Logger.Fatal(e.Start(":8080"))
}

func helloHandler(ctx echo.Context) error {
	time.Sleep(3 * time.Second)
	return ctx.String(http.StatusOK, "Hello World")
}
