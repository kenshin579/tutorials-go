// Echo HTTP 서버 + Grafana Pyroscope 연동 예제
//
// Pyroscope SDK와 Profiling Labels(TagWrapper)를 활용하여
// 엔드포인트별 프로파일링 데이터를 구분할 수 있다.
//
// 실행:
//
//	PYROSCOPE_SERVER=http://localhost:4040 go run .
//
// 부하 생성:
//
//	curl http://localhost:7080/fast
//	curl http://localhost:7080/slow
//	curl http://localhost:7080/memory
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/grafana/pyroscope-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	serverAddr := getEnv("PYROSCOPE_SERVER", "http://localhost:4040")

	profiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "echo.server",
		ServerAddress:   serverAddr,
		Logger:          pyroscope.StandardLogger,
		Tags:            map[string]string{"hostname": hostname()},
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	if err != nil {
		log.Fatalf("pyroscope 시작 실패: %v", err)
	}
	defer profiler.Stop()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/fast", handleFast)
	e.GET("/slow", handleSlow)
	e.GET("/memory", handleMemory)

	port := getEnv("PORT", "7080")
	log.Printf("서버 시작: http://localhost:%s (Pyroscope: %s)", port, serverAddr)
	e.Logger.Fatal(e.Start(":" + port))
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
