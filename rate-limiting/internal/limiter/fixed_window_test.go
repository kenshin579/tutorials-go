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

func setupMiniredis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { client.Close() })

	return mr, client
}

func TestFixedWindow_AllowWithinLimit(t *testing.T) {
	_, client := setupMiniredis(t)
	fw := NewFixedWindow(client, 5, time.Minute)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		result, err := fw.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
		assert.Equal(t, 5, result.Limit)
		assert.Equal(t, 5-i-1, result.Remaining)
	}
}

func TestFixedWindow_RejectOverLimit(t *testing.T) {
	_, client := setupMiniredis(t)
	fw := NewFixedWindow(client, 3, time.Minute)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		result, err := fw.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	result, err := fw.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, 0, result.Remaining)
}

func TestFixedWindow_ResetAfterWindow(t *testing.T) {
	mr, client := setupMiniredis(t)
	fw := NewFixedWindow(client, 2, time.Minute)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		result, err := fw.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	result, err := fw.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.False(t, result.Allowed)

	// Fast-forward time to expire the window
	mr.FastForward(time.Minute + time.Second)

	result, err = fw.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.True(t, result.Allowed)
}

func TestFixedWindow_ExactLimitBoundary(t *testing.T) {
	_, client := setupMiniredis(t)
	fw := NewFixedWindow(client, 1, time.Minute)
	ctx := context.Background()

	result, err := fw.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.Equal(t, 0, result.Remaining)

	result, err = fw.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, 0, result.Remaining)
}

func TestFixedWindow_IndependentKeys(t *testing.T) {
	_, client := setupMiniredis(t)
	fw := NewFixedWindow(client, 2, time.Minute)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		result, err := fw.Allow(ctx, "user1")
		require.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	// user1 is exhausted
	result, err := fw.Allow(ctx, "user1")
	require.NoError(t, err)
	assert.False(t, result.Allowed)

	// user2 should still be allowed
	result, err = fw.Allow(ctx, "user2")
	require.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.Equal(t, 1, result.Remaining)
}
