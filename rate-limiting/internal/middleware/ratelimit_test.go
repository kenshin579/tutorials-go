package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/kenshin579/tutorials-go/rate-limiting/internal/limiter"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEcho(t *testing.T, lim limiter.Limiter) *echo.Echo {
	t.Helper()
	e := echo.New()
	e.Use(RateLimitMiddleware(lim, IPKeyFunc))
	e.GET("/api/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
	})
	return e
}

func TestMiddleware_NormalRequest_ReturnsRateLimitHeaders(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	fw := limiter.NewFixedWindow(client, 10, time.Minute)
	e := setupTestEcho(t, fw)

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("X-Real-Ip", "192.168.1.1")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "10", rec.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "9", rec.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Reset"))
}

func TestMiddleware_ExceedLimit_Returns429(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	fw := limiter.NewFixedWindow(client, 2, time.Minute)
	e := setupTestEcho(t, fw)

	// Consume all
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("X-Real-Ip", "192.168.1.1")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 3rd request should be rejected
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("X-Real-Ip", "192.168.1.1")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.Equal(t, "0", rec.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rec.Header().Get("Retry-After"))
}

func TestMiddleware_DifferentIPs_IndependentLimits(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	fw := limiter.NewFixedWindow(client, 1, time.Minute)
	e := setupTestEcho(t, fw)

	// IP 1 - use up quota
	req1 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req1.Header.Set("X-Real-Ip", "10.0.0.1")
	rec1 := httptest.NewRecorder()
	e.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusOK, rec1.Code)

	// IP 1 - should be blocked
	req2 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req2.Header.Set("X-Real-Ip", "10.0.0.1")
	rec2 := httptest.NewRecorder()
	e.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusTooManyRequests, rec2.Code)

	// IP 2 - should still be allowed
	req3 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req3.Header.Set("X-Real-Ip", "10.0.0.2")
	rec3 := httptest.NewRecorder()
	e.ServeHTTP(rec3, req3)
	assert.Equal(t, http.StatusOK, rec3.Code)
}
