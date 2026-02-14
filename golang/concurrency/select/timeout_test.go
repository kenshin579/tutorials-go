package _select

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTimeoutWithTimeAfter - time.After로 타임아웃 처리
func TestTimeoutWithTimeAfter(t *testing.T) {
	ch := make(chan string)

	// 느린 작업 시뮬레이션
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- "result"
	}()

	var result string
	var timedOut bool

	select {
	case msg := <-ch:
		result = msg
	case <-time.After(50 * time.Millisecond):
		timedOut = true
	}

	assert.True(t, timedOut)
	assert.Empty(t, result)
}

// TestTimeoutSuccess - 시간 내에 결과를 받는 경우
func TestTimeoutSuccess(t *testing.T) {
	ch := make(chan string)

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch <- "fast result"
	}()

	var result string
	select {
	case msg := <-ch:
		result = msg
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}

	assert.Equal(t, "fast result", result)
}

// TestTimeoutWithContext - context.WithTimeout으로 타임아웃 처리
func TestTimeoutWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	ch := make(chan string)

	// 느린 작업
	go func() {
		time.Sleep(200 * time.Millisecond)
		select {
		case ch <- "result":
		case <-ctx.Done():
			return
		}
	}()

	select {
	case msg := <-ch:
		t.Fatalf("unexpected result: %s", msg)
	case <-ctx.Done():
		assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	}
}

// simulateAPICall - timeout이 있는 API 호출 시뮬레이션
func simulateAPICall(ctx context.Context, delay time.Duration) (string, error) {
	ch := make(chan string, 1)

	go func() {
		time.Sleep(delay)
		ch <- "api response"
	}()

	select {
	case result := <-ch:
		return result, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// TestSimulateAPICallSuccess - API 호출 성공
func TestSimulateAPICallSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result, err := simulateAPICall(ctx, 10*time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, "api response", result)
}

// TestSimulateAPICallTimeout - API 호출 타임아웃
func TestSimulateAPICallTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result, err := simulateAPICall(ctx, 200*time.Millisecond)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Empty(t, result)
}
