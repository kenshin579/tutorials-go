package main

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	v1 := e.Group("/api")

	v1.POST("/test", CreateEmployee)
	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}

func CreateEmployee(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Error(err)
	}
	log.Info(string(body))
	return c.NoContent(http.StatusOK)
}
