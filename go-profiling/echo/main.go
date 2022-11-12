package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echopprof "github.com/sevenNt/echo-pprof"
)

type messageRequest struct {
	Message string `json:"message"`
}

type cpuRequest struct {
	RepeatCount int `json:"repeatCount"`
}

type memoryRequest struct {
	RepeatCount int `json:"repeatCount"`
}

func main() {
	e := echo.New()
	echopprof.Wrap(e)

	e.GET("/hello", helloHandler)
	e.POST("/stress/cpu", cpuHandler)
	e.POST("/stress/memory", memoryHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func memoryHandler(ctx echo.Context) error {
	request := memoryRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	AllocMemory(request.RepeatCount)

	return ctx.NoContent(http.StatusOK)
}

func cpuHandler(ctx echo.Context) error {
	request := cpuRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	IncreaseInt(request.RepeatCount)

	return ctx.NoContent(http.StatusOK)
}

func helloHandler(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, messageRequest{
		Message: "Hello World",
	})
}

func AllocMemory(repeatCount int) {
	bytes1000 := alloc1000()
	bytes1000[0] = '0'

	for i := 0; i < repeatCount; i++ {
		time.Sleep(400 * time.Millisecond)
	}
}

func alloc1000() []byte {
	return make([]byte, 1000)
}

func IncreaseInt(repeatCount int) {
	i := 0
	for j := 0; j < repeatCount; j++ {
		i = increase1000(i)
		i = increase2000(i)
	}
	fmt.Println("cpu ended")
}

func increase1000(n int) int {
	for n := 0; n < 1000; n++ {
		n = n + 1
	}
	return n
}

func increase2000(n int) int {
	for n := 0; n < 2000; n++ {
		n = n + 1
	}
	return n
}
