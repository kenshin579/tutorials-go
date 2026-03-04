package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func setupEcho(config RateLimitConfig) *echo.Echo {
	e := echo.New()
	e.Use(RateLimitMiddleware(config))
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	return e
}

func TestMiddleware_AllowNormalRequests(t *testing.T) {
	e := setupEcho(RateLimitConfig{Rate: rate.Limit(10), Burst: 5})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddleware_RejectExcessRequests(t *testing.T) {
	e := setupEcho(RateLimitConfig{Rate: rate.Limit(1), Burst: 2})

	// Burst 2개 허용
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	}

	// 3번째 요청은 429
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestMiddleware_PerIPIsolation(t *testing.T) {
	e := setupEcho(RateLimitConfig{Rate: rate.Limit(1), Burst: 1})

	// IP1: 첫 요청 허용
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.Header.Set("X-Real-Ip", "1.1.1.1")
	rec1 := httptest.NewRecorder()
	e.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusOK, rec1.Code)

	// IP1: 두 번째 요청 거부
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.Header.Set("X-Real-Ip", "1.1.1.1")
	rec2 := httptest.NewRecorder()
	e.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusTooManyRequests, rec2.Code)

	// IP2: 다른 IP는 허용
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req3.Header.Set("X-Real-Ip", "2.2.2.2")
	rec3 := httptest.NewRecorder()
	e.ServeHTTP(rec3, req3)
	assert.Equal(t, http.StatusOK, rec3.Code)
}

func TestMiddleware_Skipper(t *testing.T) {
	e := setupEcho(RateLimitConfig{
		Rate:  rate.Limit(1),
		Burst: 1,
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/test"
		},
	})

	// Skipper가 true이면 rate limit 무시
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
