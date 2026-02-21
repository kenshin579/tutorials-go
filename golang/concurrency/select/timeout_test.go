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

	go func() {
		time.Sleep(200 * time.Millisecond) // 200ms 걸리는 느린 작업 시뮬레이션
		ch <- "result"
	}()

	var result string
	var timedOut bool

	select {
	case msg := <-ch: // 작업 결과가 먼저 오면 정상 처리
		result = msg
	case <-time.After(50 * time.Millisecond): // 50ms 초과 시 timeout channel에서 값 수신
		timedOut = true
	}

	assert.True(t, timedOut)
	assert.Empty(t, result)
}

// TestTimeoutSuccess - 시간 내에 결과를 받는 경우
func TestTimeoutSuccess(t *testing.T) {
	ch := make(chan string)

	go func() {
		time.Sleep(10 * time.Millisecond) // 10ms만에 완료되는 빠른 작업
		ch <- "fast result"
	}()

	var result string
	select {
	case msg := <-ch: // 100ms timeout 전에 결과 수신
		result = msg
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}

	assert.Equal(t, "fast result", result)
}

// TestTimeoutWithContext - context.WithTimeout으로 타임아웃 처리
func TestTimeoutWithContext(t *testing.T) {
	// 50ms timeout 설정
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel() // 리소스 해제를 위해 반드시 cancel 호출

	ch := make(chan string)

	go func() {
		time.Sleep(200 * time.Millisecond) // 200ms 걸리는 느린 작업
		select {
		case ch <- "result":
		case <-ctx.Done(): // context 취소 시 goroutine 정리
			return
		}
	}()

	select {
	case msg := <-ch:
		t.Fatalf("unexpected result: %s", msg)
	case <-ctx.Done(): // context timeout 초과 시 에러 반환
		assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	}
}

// simulateAPICall - context 기반 timeout이 적용된 API 호출 시뮬레이션
func simulateAPICall(ctx context.Context, delay time.Duration) (string, error) {
	ch := make(chan string, 1) // 버퍼 1: goroutine이 결과를 보내고 바로 종료 가능

	go func() {
		time.Sleep(delay) // API 호출 지연 시뮬레이션
		ch <- "api response"
	}()

	select {
	case result := <-ch: // API 응답이 먼저 오면 정상 반환
		return result, nil
	case <-ctx.Done(): // context timeout 초과 시 에러 반환
		return "", ctx.Err()
	}
}

// TestSimulateAPICallSuccess - API 호출 성공
func TestSimulateAPICallSuccess(t *testing.T) {
	// 100ms timeout — API는 10ms만에 완료
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result, err := simulateAPICall(ctx, 10*time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, "api response", result)
}

// TestSimulateAPICallTimeout - API 호출 타임아웃
func TestSimulateAPICallTimeout(t *testing.T) {
	// 50ms timeout — API는 200ms 걸리므로 timeout 발생
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result, err := simulateAPICall(ctx, 200*time.Millisecond)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Empty(t, result)
}
