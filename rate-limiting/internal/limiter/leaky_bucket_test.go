package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newLeakyBucketTestClient(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { client.Close() })

	return mr, client
}

func TestLeakyBucket_AllowWithinCapacity(t *testing.T) {
	_, client := newLeakyBucketTestClient(t)
	lb := NewLeakyBucket(client, 5, 1.0)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		result, err := lb.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
		assert.Equal(t, 5, result.Limit)
	}
}

func TestLeakyBucket_RejectWhenFull(t *testing.T) {
	_, client := newLeakyBucketTestClient(t)
	lb := NewLeakyBucket(client, 3, 1.0)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		result, err := lb.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	result, err := lb.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, 0, result.Remaining)
}

func TestLeakyBucket_LeakOverTime(t *testing.T) {
	_, client := newLeakyBucketTestClient(t)
	clock, advance := fakeClock(time.Now())
	lb := NewLeakyBucketWithClock(client, 3, 1.0, clock)
	ctx := context.Background()

	// Fill the bucket
	for i := 0; i < 3; i++ {
		result, err := lb.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	// Queue is full
	result, err := lb.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.False(t, result.Allowed)

	// Advance clock by 2 seconds -> 2 items leaked
	advance(2 * time.Second)

	result, err = lb.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.True(t, result.Allowed)
}
