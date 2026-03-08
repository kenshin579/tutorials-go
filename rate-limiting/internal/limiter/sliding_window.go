package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var slidingWindowScript = redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local current_window = math.floor(now / window) * window
local prev_window = current_window - window

local prev_key = key .. ":" .. prev_window
local curr_key = key .. ":" .. current_window

local prev_count = tonumber(redis.call("GET", prev_key) or "0") or 0
local curr_count = tonumber(redis.call("GET", curr_key) or "0") or 0

local elapsed = now - current_window
local weight = 1 - (elapsed / window)
local weighted_count = math.floor(prev_count * weight + curr_count)

if weighted_count < limit then
    redis.call("INCR", curr_key)
    redis.call("EXPIRE", curr_key, window * 2)
    curr_count = curr_count + 1
    weighted_count = math.floor(prev_count * weight + curr_count)
    local remaining = limit - weighted_count
    if remaining < 0 then remaining = 0 end
    local reset_at = current_window + window
    return {1, remaining, reset_at}
else
    local remaining = 0
    local reset_at = current_window + window
    return {0, remaining, reset_at}
end
`)

// SlidingWindow implements the Sliding Window Counter algorithm.
// It uses a weighted average of the previous and current window counts.
type SlidingWindow struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewSlidingWindow(client *redis.Client, limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (sw *SlidingWindow) Allow(ctx context.Context, key string) (*Result, error) {
	windowSec := int(sw.window.Seconds())
	now := time.Now()
	nowSec := now.Unix()

	redisKey := fmt.Sprintf("ratelimit:sw:%s", key)

	res, err := slidingWindowScript.Run(ctx, sw.client, []string{redisKey}, sw.limit, windowSec, nowSec).Int64Slice()
	if err != nil {
		return nil, fmt.Errorf("sliding window script error: %w", err)
	}

	allowed := res[0] == 1
	remaining := int(res[1])
	resetAt := time.Unix(res[2], 0)

	return &Result{
		Allowed:   allowed,
		Limit:     sw.limit,
		Remaining: remaining,
		ResetAt:   resetAt,
	}, nil
}
