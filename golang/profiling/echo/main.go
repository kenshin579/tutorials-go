// Echo 프레임워크에서 pprof 프로파일링을 사용하는 예제
//
// echo-pprof 라이브러리를 사용하면 Echo 라우터에 pprof 엔드포인트를 쉽게 추가할 수 있다.
// 기본 net/http/pprof는 DefaultServeMux에 등록되지만, Echo는 자체 라우터를 사용하므로
// echopprof.Wrap()으로 별도 등록이 필요하다.
//
// 실행 후 프로파일링:
//
//	curl -X POST http://localhost:8080/stress/cpu -H "Content-Type: application/json" -d '{"repeatCount": 100}'
//	go tool pprof http://localhost:8080/debug/pprof/profile?seconds=10
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
	echopprof.Wrap(e) // Echo 라우터에 /debug/pprof/* 엔드포인트 등록

	e.GET("/hello", helloHandler)
	e.POST("/stress/cpu", cpuHandler)       // CPU 부하 생성 API
	e.POST("/stress/memory", memoryHandler) // 메모리 부하 생성 API

	e.Logger.Fatal(e.Start(":8080"))
}

// memoryHandler는 요청 횟수만큼 메모리를 할당하고 대기하여 힙 프로파일에서 확인할 수 있도록 한다.
func memoryHandler(ctx echo.Context) error {
	request := memoryRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	AllocMemory(request.RepeatCount)

	return ctx.NoContent(http.StatusOK)
}

// cpuHandler는 요청 횟수만큼 CPU 연산을 수행하여 CPU 프로파일에서 확인할 수 있도록 한다.
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

// AllocMemory는 1000바이트를 힙에 할당한 후 repeatCount만큼 대기한다.
// 할당 후 대기하므로 pprof inuse_space에서 메모리가 잡힌 상태를 확인할 수 있다.
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

// IncreaseInt는 repeatCount만큼 CPU 부하를 생성한다.
// increase1000과 increase2000의 CPU 사용 비율을 프로파일에서 비교할 수 있다.
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
