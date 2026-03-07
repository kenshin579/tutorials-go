package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"

	"grafana-tracing/metrics"
)

// PrometheusMiddleware는 모든 HTTP 요청에 대해 메트릭을 자동 수집한다 (Exemplar 포함)
func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/metrics" {
				return next(c)
			}

			start := time.Now()

			metrics.HttpRequestsInFlight.Inc()
			defer metrics.HttpRequestsInFlight.Dec()

			err := next(c)

			status := c.Response().Status
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				}
			}

			duration := time.Since(start).Seconds()
			method := c.Request().Method
			path := c.Path()

			metrics.HttpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()

			// Exemplar가 포함된 Histogram 기록
			observeWithExemplar(metrics.HttpRequestDuration, c.Request().Context(), duration, method, path)

			return err
		}
	}
}

// observeWithExemplar는 Trace ID를 Exemplar로 첨부해서 Histogram에 기록한다
func observeWithExemplar(histogram *prometheus.HistogramVec, ctx context.Context, duration float64, labels ...string) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		histogram.WithLabelValues(labels...).(prometheus.ExemplarObserver).ObserveWithExemplar(
			duration, prometheus.Labels{
				"traceID": spanCtx.TraceID().String(),
			},
		)
	} else {
		histogram.WithLabelValues(labels...).Observe(duration)
	}
}
