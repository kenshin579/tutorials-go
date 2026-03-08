package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var tokenBucketScript = redis.NewScript(`
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local data = redis.call("HMGET", key, "tokens", "last_refill")
local tokens = tonumber(data[1])
local last_refill = tonumber(data[2])

if tokens == nil then
    tokens = capacity
    last_refill = now
end

local elapsed = now - last_refill
local new_tokens = elapsed * refill_rate
tokens = math.min(capacity, tokens + new_tokens)

if tokens >= 1 then
    tokens = tokens - 1
    redis.call("HMSET", key, "tokens", tokens, "last_refill", now)
    redis.call("EXPIRE", key, math.ceil(capacity / refill_rate) * 2)
    return {1, math.floor(tokens)}
else
    redis.call("HMSET", key, "tokens", tokens, "last_refill", now)
    redis.call("EXPIRE", key, math.ceil(capacity / refill_rate) * 2)
    return {0, 0}
end
`)

// TokenBucket implements the Token Bucket algorithm.
// Tokens are added at a fixed rate up to a maximum capacity.
// Each request consumes one token; requests are rejected when no tokens remain.
type TokenBucket struct {
	client     *redis.Client
	capacity   int
	refillRate float64 // tokens per second
	clock      Clock
}

func NewTokenBucket(client *redis.Client, capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		client:     client,
		capacity:   capacity,
		refillRate: refillRate,
		clock:      RealClock,
	}
}

func NewTokenBucketWithClock(client *redis.Client, capacity int, refillRate float64, clock Clock) *TokenBucket {
	return &TokenBucket{
		client:     client,
		capacity:   capacity,
		refillRate: refillRate,
		clock:      clock,
	}
}

func (tb *TokenBucket) Allow(ctx context.Context, key string) (*Result, error) {
	now := tb.clock()
	nowSec := float64(now.UnixMilli()) / 1000.0
	redisKey := fmt.Sprintf("ratelimit:tb:%s", key)

	res, err := tokenBucketScript.Run(ctx, tb.client, []string{redisKey}, tb.capacity, tb.refillRate, nowSec).Int64Slice()
	if err != nil {
		return nil, fmt.Errorf("token bucket script error: %w", err)
	}

	allowed := res[0] == 1
	remaining := int(res[1])

	// Calculate when one token will be available
	refillTime := time.Duration(float64(time.Second) / tb.refillRate)
	resetAt := now.Add(refillTime)

	return &Result{
		Allowed:   allowed,
		Limit:     tb.capacity,
		Remaining: remaining,
		ResetAt:   resetAt,
	}, nil
}
