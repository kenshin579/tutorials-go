package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	e := echo.New()

	// ============================================================
	// 로깅/추적
	// ============================================================

	// Logger - 요청/응답 로깅 (RequestLoggerWithConfig 사용)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("method", values.Method),
				zap.String("uri", values.URI),
				zap.Int("status", values.Status),
			)
			return nil
		},
	}))

	// RequestID - 요청 추적 ID 자동 생성
	e.Use(middleware.RequestID())

	// ============================================================
	// 보안
	// ============================================================

	// CORS - 교차 출처 리소스 공유
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://example.com", "https://app.example.com"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Secure - 보안 헤더 (XSS, HSTS 등)
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	// ============================================================
	// 안정성/성능
	// ============================================================

	// Recover - 패닉 복구
	e.Use(middleware.Recover())

	// RateLimiter - 요청 속도 제한 (초당 20회)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// BodyLimit - 요청 본문 크기 제한 (2MB)
	e.Use(middleware.BodyLimit("2M"))

	// Gzip - 응답 압축
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level:     5,
		MinLength: 1024, // 1KB 이상만 압축
	}))

	// Timeout - 요청 타임아웃
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	// ============================================================
	// 라우트
	// ============================================================

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message":    "Echo 빌트인 미들웨어 예제",
			"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
