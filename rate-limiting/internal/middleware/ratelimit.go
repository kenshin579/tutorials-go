package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kenshin579/tutorials-go/rate-limiting/internal/limiter"
	"github.com/labstack/echo/v4"
)

// KeyFunc extracts a rate limit key from the request.
type KeyFunc func(c echo.Context) string

// IPKeyFunc returns the client IP as the rate limit key.
func IPKeyFunc(c echo.Context) string {
	return c.RealIP()
}

// RateLimitMiddleware creates an Echo middleware that enforces rate limiting.
func RateLimitMiddleware(lim limiter.Limiter, keyFunc KeyFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := keyFunc(c)

			result, err := lim.Allow(c.Request().Context(), key)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "rate limiter error")
			}

			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(result.Limit))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))

			if !result.Allowed {
				retryAfter := time.Until(result.ResetAt).Seconds()
				if retryAfter < 1 {
					retryAfter = 1
				}
				c.Response().Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter))
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "rate limit exceeded",
				})
			}

			return next(c)
		}
	}
}
