package main

import (
	"context"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/labstack/echo/v4"
)

type response struct {
	Message string `json:"message"`
	Elapsed string `json:"elapsed,omitempty"`
}

// handleFast는 단순 JSON 응답을 반환한다 (비교 기준선).
func handleFast(c echo.Context) error {
	return c.JSON(http.StatusOK, response{Message: "fast response"})
}

// handleSlow는 CPU 집약적 연산을 수행한다.
// Pyroscope Labels로 태깅하여 Flame Graph에서 /slow 엔드포인트만 필터링할 수 있다.
func handleSlow(c echo.Context) error {
	start := time.Now()

	pyroscope.TagWrapper(c.Request().Context(),
		pyroscope.Labels("endpoint", "/slow"),
		func(ctx context.Context) {
			// 피보나치 재귀 연산으로 CPU 부하
			fibonacci(38)
		})

	return c.JSON(http.StatusOK, response{
		Message: "slow response (CPU intensive)",
		Elapsed: time.Since(start).String(),
	})
}

// handleMemory는 대량 메모리를 할당한다.
// Pyroscope Labels로 태깅하여 Flame Graph에서 /memory 엔드포인트만 필터링할 수 있다.
func handleMemory(c echo.Context) error {
	start := time.Now()

	pyroscope.TagWrapper(c.Request().Context(),
		pyroscope.Labels("endpoint", "/memory"),
		func(ctx context.Context) {
			allocateMemory()
		})

	return c.JSON(http.StatusOK, response{
		Message: "memory response (heap allocation)",
		Elapsed: time.Since(start).String(),
	})
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// allocateMemory는 여러 고루틴에서 동시에 메모리를 할당한다.
func allocateMemory() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := make([]byte, 5*1024*1024) // 5MB
			for j := range data {
				data[j] = byte(j % 256)
			}
			_ = math.Sqrt(float64(len(data)))
			time.Sleep(50 * time.Millisecond)
		}()
	}
	wg.Wait()
}
