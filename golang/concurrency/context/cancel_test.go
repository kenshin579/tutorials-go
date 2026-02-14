package context_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestWithCancel - context.WithCancel 기본 사용법
func TestWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var stopped atomic.Bool

	go func() {
		<-ctx.Done() // cancel() 호출까지 대기
		stopped.Store(true)
	}()

	cancel() // 명시적으로 취소
	time.Sleep(50 * time.Millisecond)

	assert.True(t, stopped.Load())
	assert.ErrorIs(t, ctx.Err(), context.Canceled)
}

// TestCancelChain - parent context 취소 시 child도 취소
func TestCancelChain(t *testing.T) {
	parent, parentCancel := context.WithCancel(context.Background())
	child, childCancel := context.WithCancel(parent)
	defer childCancel()

	parentCancel() // parent 취소 → child도 취소
	time.Sleep(50 * time.Millisecond)

	assert.ErrorIs(t, parent.Err(), context.Canceled)
	assert.ErrorIs(t, child.Err(), context.Canceled)
}

// TestCancelWorker - cancel로 goroutine 정리
func TestCancelWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	results := make(chan int, 10)

	// worker: cancel될 때까지 작업
	go func() {
		defer close(results)
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			case results <- i:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// 50ms 후 cancel
	time.Sleep(50 * time.Millisecond)
	cancel()

	var collected []int
	for v := range results {
		collected = append(collected, v)
	}

	t.Logf("수집된 값: %v", collected)
	assert.Greater(t, len(collected), 0)
}
