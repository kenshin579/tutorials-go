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

func newSlidingWindowTestClient(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { client.Close() })

	return mr, client
}

func TestSlidingWindow_AllowWithinLimit(t *testing.T) {
	_, client := newSlidingWindowTestClient(t)
	sw := NewSlidingWindow(client, 5, 10*time.Second)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		result, err := sw.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}
}

func TestSlidingWindow_RejectOverLimit(t *testing.T) {
	_, client := newSlidingWindowTestClient(t)
	sw := NewSlidingWindow(client, 3, 10*time.Second)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		result, err := sw.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	result, err := sw.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, 0, result.Remaining)
}

func TestSlidingWindow_ResetAfterWindowExpiry(t *testing.T) {
	mr, client := newSlidingWindowTestClient(t)
	sw := NewSlidingWindow(client, 2, 10*time.Second)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		result, err := sw.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	result, err := sw.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.False(t, result.Allowed)

	// Fast-forward past two full windows so previous count fully expires
	mr.FastForward(21 * time.Second)

	result, err = sw.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.True(t, result.Allowed)
}
