package ratelimit

import (
	"context"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

// DistributedRateLimiter provides Redis-based distributed rate limiting using GCRA algorithm.
type DistributedRateLimiter struct {
	limiter *redis_rate.Limiter
}

// NewDistributedRateLimiter creates a new distributed rate limiter backed by Redis.
func NewDistributedRateLimiter(rdb *redis.Client) *DistributedRateLimiter {
	return &DistributedRateLimiter{
		limiter: redis_rate.NewLimiter(rdb),
	}
}

// Allow checks whether a request identified by key is allowed under the given limit.
func (d *DistributedRateLimiter) Allow(ctx context.Context, key string, limit redis_rate.Limit) (*redis_rate.Result, error) {
	return d.limiter.Allow(ctx, key, limit)
}
