package main

import (
	"log"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"grafana-metrics/handler"
	"grafana-metrics/middleware"
)

func main() {
	e := echo.New()

	// 기본 미들웨어
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())

	// Prometheus 메트릭 자동 수집 미들웨어
	e.Use(middleware.PrometheusMiddleware())

	// /metrics 엔드포인트
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// 헬스체크
	e.GET("/health", handler.HealthCheck)

	// 주문 API
	e.POST("/api/orders", handler.CreateOrder)
	e.GET("/api/orders", handler.ListOrders)
	e.GET("/api/orders/:id", handler.GetOrder)

	log.Println("Server starting on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
