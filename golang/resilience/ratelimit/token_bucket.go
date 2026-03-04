package ratelimit

import (
	"context"

	"golang.org/x/time/rate"
)

// RateLimiter wraps x/time/rate.Limiter providing three usage patterns.
type RateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter creates a new RateLimiter.
// r is the rate of tokens per second, burst is the maximum burst size.
func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{limiter: rate.NewLimiter(r, burst)}
}

// Allow reports whether an event may happen now (non-blocking).
func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}

// Wait blocks until the limiter permits an event or ctx is done.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}

// Reserve returns a Reservation that indicates how long the caller must wait.
func (rl *RateLimiter) Reserve() *rate.Reservation {
	return rl.limiter.Reserve()
}
