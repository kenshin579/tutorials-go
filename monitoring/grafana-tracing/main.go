package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"grafana-tracing/handler"
	"grafana-tracing/middleware"
	"grafana-tracing/tracing"
)

func main() {
	ctx := context.Background()

	// OpenTelemetry TracerProvider 초기화
	tp, err := tracing.InitTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	e := echo.New()

	// 기본 미들웨어
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())

	// OTel Echo 미들웨어 (HTTP 요청 자동 트레이싱)
	e.Use(otelecho.Middleware("order-service"))

	// Prometheus 메트릭 미들웨어 (Exemplar 포함)
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
