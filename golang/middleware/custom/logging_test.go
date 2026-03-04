package custom

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLogger(t *testing.T) {
	// 로그 캡처를 위한 observer 설정
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	e := echo.New()
	e.Use(ZapLogger(logger))
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, recorded.Len())

	entry := recorded.All()[0]
	assert.Equal(t, "request", entry.Message)
	assert.Equal(t, "GET", entry.ContextMap()["method"])
	assert.Equal(t, "/test", entry.ContextMap()["path"])
	assert.Equal(t, int64(200), entry.ContextMap()["status"])
}

func TestZapLogger_ErrorRequest(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	e := echo.New()
	e.Use(ZapLogger(logger))
	e.GET("/error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, 1, recorded.Len())

	entry := recorded.All()[0]
	assert.Equal(t, "request", entry.Message)
	assert.Equal(t, zap.ErrorLevel, entry.Level)
}

func TestZapLogger_Skipper(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	e := echo.New()
	e.Use(ZapLoggerWithConfig(logger, ZapLoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/health"
		},
	}))
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 0, recorded.Len()) // Skipper로 건너뛰어 로그 없음
}
