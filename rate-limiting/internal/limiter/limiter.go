package limiter

import (
	"context"
	"time"
)

// Result contains the rate limiting decision and metadata for HTTP headers.
type Result struct {
	Allowed   bool
	Limit     int
	Remaining int
	ResetAt   time.Time
}

// Limiter is the interface that all rate limiting algorithms must implement.
type Limiter interface {
	Allow(ctx context.Context, key string) (*Result, error)
}

// Clock provides the current time. Override for testing.
type Clock func() time.Time

// RealClock returns time.Now.
func RealClock() time.Time {
	return time.Now()
}
