package ratelimit

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// RateLimitConfig defines the config for rate limit middleware.
type RateLimitConfig struct {
	Rate    rate.Limit
	Burst   int
	Skipper middleware.Skipper
}

// RateLimitMiddleware returns Echo middleware that limits requests per IP.
func RateLimitMiddleware(config RateLimitConfig) echo.MiddlewareFunc {
	var clients sync.Map

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}

			ip := c.RealIP()
			v, _ := clients.LoadOrStore(ip, rate.NewLimiter(config.Rate, config.Burst))
			limiter := v.(*rate.Limiter)

			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "rate limit exceeded",
				})
			}

			return next(c)
		}
	}
}
