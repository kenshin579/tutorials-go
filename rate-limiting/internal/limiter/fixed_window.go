package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var fixedWindowScript = redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])

local current = redis.call("INCR", key)
if current == 1 then
    redis.call("EXPIRE", key, window)
end

local ttl = redis.call("TTL", key)
if ttl < 0 then
    ttl = window
end

return {current, ttl}
`)

// FixedWindow implements the Fixed Window Counter algorithm.
// It counts requests within a fixed time window and rejects when the limit is exceeded.
type FixedWindow struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewFixedWindow(client *redis.Client, limit int, window time.Duration) *FixedWindow {
	return &FixedWindow{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (fw *FixedWindow) Allow(ctx context.Context, key string) (*Result, error) {
	windowSec := int(fw.window.Seconds())
	now := time.Now()
	windowStart := now.Truncate(fw.window)
	redisKey := fmt.Sprintf("ratelimit:fw:%s:%d", key, windowStart.Unix())

	res, err := fixedWindowScript.Run(ctx, fw.client, []string{redisKey}, fw.limit, windowSec).Int64Slice()
	if err != nil {
		return nil, fmt.Errorf("fixed window script error: %w", err)
	}

	current := int(res[0])
	ttl := time.Duration(res[1]) * time.Second

	remaining := fw.limit - current
	if remaining < 0 {
		remaining = 0
	}

	return &Result{
		Allowed:   current <= fw.limit,
		Limit:     fw.limit,
		Remaining: remaining,
		ResetAt:   now.Add(ttl),
	}, nil
}
