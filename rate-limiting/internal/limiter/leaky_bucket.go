package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var leakyBucketScript = redis.NewScript(`
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local leak_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local data = redis.call("HMGET", key, "queue_size", "last_leak")
local queue_size = tonumber(data[1]) or 0
local last_leak = tonumber(data[2]) or now

local elapsed = now - last_leak
local leaked = math.floor(elapsed * leak_rate)
queue_size = math.max(0, queue_size - leaked)

if leaked > 0 then
    last_leak = now
end

if queue_size < capacity then
    queue_size = queue_size + 1
    redis.call("HMSET", key, "queue_size", queue_size, "last_leak", last_leak)
    redis.call("EXPIRE", key, math.ceil(capacity / leak_rate) * 2)
    local remaining = capacity - queue_size
    return {1, remaining}
else
    redis.call("HMSET", key, "queue_size", queue_size, "last_leak", last_leak)
    redis.call("EXPIRE", key, math.ceil(capacity / leak_rate) * 2)
    return {0, 0}
end
`)

// LeakyBucket implements the Leaky Bucket algorithm.
// Requests fill a queue that leaks at a constant rate.
// When the queue is full, new requests are rejected.
type LeakyBucket struct {
	client   *redis.Client
	capacity int
	leakRate float64 // requests leaked per second
	clock    Clock
}

func NewLeakyBucket(client *redis.Client, capacity int, leakRate float64) *LeakyBucket {
	return &LeakyBucket{
		client:   client,
		capacity: capacity,
		leakRate: leakRate,
		clock:    RealClock,
	}
}

func NewLeakyBucketWithClock(client *redis.Client, capacity int, leakRate float64, clock Clock) *LeakyBucket {
	return &LeakyBucket{
		client:   client,
		capacity: capacity,
		leakRate: leakRate,
		clock:    clock,
	}
}

func (lb *LeakyBucket) Allow(ctx context.Context, key string) (*Result, error) {
	now := lb.clock()
	nowSec := float64(now.UnixMilli()) / 1000.0
	redisKey := fmt.Sprintf("ratelimit:lb:%s", key)

	res, err := leakyBucketScript.Run(ctx, lb.client, []string{redisKey}, lb.capacity, lb.leakRate, nowSec).Int64Slice()
	if err != nil {
		return nil, fmt.Errorf("leaky bucket script error: %w", err)
	}

	allowed := res[0] == 1
	remaining := int(res[1])

	// Next leak happens in 1/leakRate seconds
	leakInterval := time.Duration(float64(time.Second) / lb.leakRate)
	resetAt := now.Add(leakInterval)

	return &Result{
		Allowed:   allowed,
		Limit:     lb.capacity,
		Remaining: remaining,
		ResetAt:   resetAt,
	}, nil
}
