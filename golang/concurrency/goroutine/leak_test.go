package goroutine

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGoroutineLeak - channel 대기로 인한 goroutine leak 예시
func TestGoroutineLeak(t *testing.T) {
	initialCount := runtime.NumGoroutine()

	// leak을 발생시키는 함수: 아무도 receive하지 않는 channel에 send 대기
	leakyFunc := func() <-chan int {
		ch := make(chan int)
		go func() {
			// 이 goroutine은 영원히 blocking된다 (누구도 ch에서 receive하지 않으므로)
			ch <- 42
		}()
		return ch
	}

	// channel을 반환받지만 사용하지 않음 → goroutine leak
	_ = leakyFunc()

	time.Sleep(50 * time.Millisecond)

	leakedCount := runtime.NumGoroutine()
	t.Logf("초기: %d, leak 후: %d", initialCount, leakedCount)

	// goroutine이 증가했음 (leak)
	assert.Greater(t, leakedCount, initialCount)
}

// TestGoroutineLeakPrevention_WithContext - context로 goroutine leak 방지
func TestGoroutineLeakPrevention_WithContext(t *testing.T) {
	initialCount := runtime.NumGoroutine()

	safeFunc := func(ctx context.Context) <-chan int {
		ch := make(chan int, 1) // buffered channel 사용
		go func() {
			defer close(ch)
			select {
			case ch <- 42:
				// 정상적으로 값 전달
			case <-ctx.Done():
				// context 취소 시 goroutine 종료
				return
			}
		}()
		return ch
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch := safeFunc(ctx)
	_ = ch

	// context를 취소하면 goroutine이 정리된다
	cancel()
	time.Sleep(50 * time.Millisecond)

	finalCount := runtime.NumGoroutine()
	t.Logf("초기: %d, 정리 후: %d", initialCount, finalCount)

	assert.LessOrEqual(t, finalCount, initialCount+1)
}

// TestGoroutineLeakPrevention_WithDone - done channel로 goroutine leak 방지
func TestGoroutineLeakPrevention_WithDone(t *testing.T) {
	initialCount := runtime.NumGoroutine()

	safeFunc := func(done <-chan struct{}) <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			select {
			case ch <- 42:
			case <-done:
				return
			}
		}()
		return ch
	}

	done := make(chan struct{})
	ch := safeFunc(done)
	_ = ch

	// done channel을 닫으면 goroutine이 정리된다
	close(done)
	time.Sleep(50 * time.Millisecond)

	finalCount := runtime.NumGoroutine()
	t.Logf("초기: %d, 정리 후: %d", initialCount, finalCount)

	assert.LessOrEqual(t, finalCount, initialCount+1)
}

// TestDetectGoroutineLeak - goroutine leak을 탐지하는 패턴
func TestDetectGoroutineLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	// 테스트 대상 코드 실행
	ch := make(chan int, 1)
	go func() {
		ch <- 1
	}()

	// 결과를 반드시 소비
	<-ch

	time.Sleep(50 * time.Millisecond)
	after := runtime.NumGoroutine()

	// goroutine 수가 크게 증가하지 않았는지 확인
	leaked := after - before
	t.Logf("goroutine 변화: %d → %d (차이: %d)", before, after, leaked)

	assert.LessOrEqual(t, leaked, 1, "goroutine leak이 발생하지 않아야 함")
}
