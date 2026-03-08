package limiter

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTokenBucketTestClient(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { client.Close() })

	return mr, client
}

// fakeClock returns a Clock that can be advanced manually.
func fakeClock(start time.Time) (Clock, func(d time.Duration)) {
	var mu sync.Mutex
	now := start
	clock := func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return now
	}
	advance := func(d time.Duration) {
		mu.Lock()
		defer mu.Unlock()
		now = now.Add(d)
	}
	return clock, advance
}

func TestTokenBucket_ConsumeTokens(t *testing.T) {
	_, client := newTokenBucketTestClient(t)
	tb := NewTokenBucket(client, 5, 1.0)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		result, err := tb.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
		assert.Equal(t, 5, result.Limit)
		assert.Equal(t, 5-i-1, result.Remaining)
	}
}

func TestTokenBucket_BurstAllowed(t *testing.T) {
	_, client := newTokenBucketTestClient(t)
	tb := NewTokenBucket(client, 10, 1.0)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		result, err := tb.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}
}

func TestTokenBucket_RejectWhenEmpty(t *testing.T) {
	_, client := newTokenBucketTestClient(t)
	tb := NewTokenBucket(client, 3, 1.0)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		result, err := tb.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	result, err := tb.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, 0, result.Remaining)
}

func TestTokenBucket_RefillAfterTime(t *testing.T) {
	_, client := newTokenBucketTestClient(t)
	clock, advance := fakeClock(time.Now())
	tb := NewTokenBucketWithClock(client, 2, 1.0, clock)
	ctx := context.Background()

	// Consume all tokens
	for i := 0; i < 2; i++ {
		result, err := tb.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	result, err := tb.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.False(t, result.Allowed)

	// Advance clock by 2 seconds -> 2 tokens refilled
	advance(2 * time.Second)

	result, err = tb.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.True(t, result.Allowed)
}
