package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedis(t *testing.T) (*redis.Client, func()) {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	endpoint, err := container.Endpoint(ctx, "")
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{Addr: endpoint})
	require.NoError(t, rdb.Ping(ctx).Err())

	cleanup := func() {
		rdb.Close()
		container.Terminate(ctx)
	}
	return rdb, cleanup
}

func TestDistributedRateLimiter_AllowAndDeny(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	rdb, cleanup := setupRedis(t)
	defer cleanup()

	ctx := context.Background()
	limiter := NewDistributedRateLimiter(rdb)
	limit := redis_rate.PerSecond(2)

	// 첫 2개 요청 허용
	for i := 0; i < 2; i++ {
		result, err := limiter.Allow(ctx, "test:user1", limit)
		require.NoError(t, err)
		assert.True(t, result.Allowed > 0, "request %d should be allowed", i+1)
	}

	// 초과 요청 거부
	result, err := limiter.Allow(ctx, "test:user1", limit)
	require.NoError(t, err)
	assert.Equal(t, 0, result.Allowed, "excess request should be denied")
	assert.True(t, result.RetryAfter > 0, "should have retry-after")
}

func TestDistributedRateLimiter_DifferentKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	rdb, cleanup := setupRedis(t)
	defer cleanup()

	ctx := context.Background()
	limiter := NewDistributedRateLimiter(rdb)
	limit := redis_rate.PerSecond(1)

	// user1 요청
	r1, err := limiter.Allow(ctx, "test:user1", limit)
	require.NoError(t, err)
	assert.True(t, r1.Allowed > 0)

	// user1 초과
	r2, err := limiter.Allow(ctx, "test:user1", limit)
	require.NoError(t, err)
	assert.Equal(t, 0, r2.Allowed)

	// user2는 별도 키이므로 허용
	r3, err := limiter.Allow(ctx, "test:user2", limit)
	require.NoError(t, err)
	assert.True(t, r3.Allowed > 0)
}

func TestDistributedRateLimiter_MultipleClients(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	rdb, cleanup := setupRedis(t)
	defer cleanup()

	ctx := context.Background()
	limiter1 := NewDistributedRateLimiter(rdb)
	limiter2 := NewDistributedRateLimiter(rdb)

	limit := redis_rate.Limit{
		Rate:   2,
		Burst:  2,
		Period: time.Second,
	}

	// 클라이언트1에서 2개 요청
	for i := 0; i < 2; i++ {
		r, err := limiter1.Allow(ctx, "shared:key", limit)
		require.NoError(t, err)
		assert.True(t, r.Allowed > 0)
	}

	// 클라이언트2에서 같은 키 요청 → 거부 (공유 카운터)
	r, err := limiter2.Allow(ctx, "shared:key", limit)
	require.NoError(t, err)
	assert.Equal(t, 0, r.Allowed, "shared counter should reject from second client")
}
