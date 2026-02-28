package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"grafana-metrics/metrics"
)

// PrometheusMiddleware는 모든 HTTP 요청에 대해 메트릭을 자동 수집한다
func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// /metrics 엔드포인트는 계측 제외
			if c.Path() == "/metrics" {
				return next(c)
			}

			start := time.Now()

			// 활성 요청 수 증가
			metrics.HttpRequestsInFlight.Inc()
			defer metrics.HttpRequestsInFlight.Dec()

			// 다음 핸들러 실행
			err := next(c)

			// 응답 정보 수집
			status := c.Response().Status
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				}
			}

			duration := time.Since(start).Seconds()
			method := c.Request().Method
			path := c.Path() // 패턴 경로 사용 (/api/orders/:id)

			// 메트릭 기록
			metrics.HttpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
			metrics.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)

			return err
		}
	}
}
