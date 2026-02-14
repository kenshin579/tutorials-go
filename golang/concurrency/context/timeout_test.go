package context_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestWithTimeout - context.WithTimeout 기본 사용법
func TestWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	select {
	case <-time.After(200 * time.Millisecond):
		t.Fatal("should have timed out")
	case <-ctx.Done():
		assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	}
}

// TestWithDeadline - context.WithDeadline 사용법
func TestWithDeadline(t *testing.T) {
	deadline := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	<-ctx.Done()
	assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)

	// Deadline() 메서드로 만료 시간 확인
	dl, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.Equal(t, deadline, dl)
}

// TestTimeoutBeforeCompletion - 작업이 timeout 전에 완료되는 경우
func TestTimeoutBeforeCompletion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ch := make(chan string, 1)
	go func() {
		time.Sleep(10 * time.Millisecond) // 빠른 작업
		ch <- "done"
	}()

	select {
	case result := <-ch:
		assert.Equal(t, "done", result)
	case <-ctx.Done():
		t.Fatal("should not timeout")
	}
}

// TestNestedTimeout - 중첩 timeout (더 짧은 것이 우선)
func TestNestedTimeout(t *testing.T) {
	outer, outerCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer outerCancel()

	inner, innerCancel := context.WithTimeout(outer, 50*time.Millisecond)
	defer innerCancel()

	<-inner.Done()

	// inner가 먼저 timeout
	assert.ErrorIs(t, inner.Err(), context.DeadlineExceeded)
	// outer는 아직 살아있음
	assert.NoError(t, outer.Err())
}
