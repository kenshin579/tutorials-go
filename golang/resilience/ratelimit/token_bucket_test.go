package ratelimit

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestRateLimiter_Allow(t *testing.T) {
	// 1 token/sec, burst 3
	rl := NewRateLimiter(rate.Limit(1), 3)

	// Burst 3개까지는 즉시 허용
	assert.True(t, rl.Allow())
	assert.True(t, rl.Allow())
	assert.True(t, rl.Allow())

	// Burst 소진 후 거부
	assert.False(t, rl.Allow())
}

func TestRateLimiter_Wait(t *testing.T) {
	// 10 tokens/sec, burst 1
	rl := NewRateLimiter(rate.Limit(10), 1)

	// 첫 번째는 즉시 허용
	ctx := context.Background()
	require.NoError(t, rl.Wait(ctx))

	// 두 번째는 대기 후 허용 (약 100ms)
	start := time.Now()
	require.NoError(t, rl.Wait(ctx))
	elapsed := time.Since(start)
	assert.True(t, elapsed >= 50*time.Millisecond, "should wait for token replenishment")
}

func TestRateLimiter_Wait_ContextCancel(t *testing.T) {
	// 매우 느린 rate
	rl := NewRateLimiter(rate.Limit(0.1), 1)

	// Burst 소진
	rl.Allow()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := rl.Wait(ctx)
	assert.Error(t, err, "should fail when context is cancelled")
}

func TestRateLimiter_Reserve(t *testing.T) {
	rl := NewRateLimiter(rate.Limit(10), 1)

	// 첫 번째 예약은 즉시
	r1 := rl.Reserve()
	assert.True(t, r1.OK())
	assert.Zero(t, r1.Delay())

	// 두 번째 예약은 대기 필요
	r2 := rl.Reserve()
	assert.True(t, r2.OK())
	assert.True(t, r2.Delay() > 0)
}

func TestRateLimiter_Concurrency(t *testing.T) {
	// 100 tokens/sec, burst 10
	rl := NewRateLimiter(rate.Limit(100), 10)

	var allowed atomic.Int64
	done := make(chan struct{})

	for i := 0; i < 50; i++ {
		go func() {
			if rl.Allow() {
				allowed.Add(1)
			}
			done <- struct{}{}
		}()
	}

	for i := 0; i < 50; i++ {
		<-done
	}

	// Burst 10이므로 최대 10개 허용
	assert.True(t, allowed.Load() <= 10, "should not exceed burst size")
	assert.True(t, allowed.Load() > 0, "should allow some requests")
}
