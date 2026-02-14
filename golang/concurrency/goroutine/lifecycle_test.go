package goroutine

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMainExitKillsGoroutines - main goroutine이 종료되면 다른 goroutine도 종료된다
// 이 테스트는 goroutine이 완료되기 전에 함수가 반환되면 작업이 유실되는 문제를 보여준다
func TestMainExitKillsGoroutines(t *testing.T) {
	var completed atomic.Bool

	go func() {
		time.Sleep(100 * time.Millisecond) // 시간이 걸리는 작업 시뮬레이션
		completed.Store(true)
	}()

	// 기다리지 않으면 goroutine은 완료되지 않는다
	// (테스트에서는 time.Sleep으로 시뮬레이션)
	assert.False(t, completed.Load(), "goroutine이 아직 완료되지 않았어야 함")
}

// TestWaitGroupSolution - sync.WaitGroup으로 goroutine 완료 대기
func TestWaitGroupSolution(t *testing.T) {
	var completed atomic.Bool
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		completed.Store(true)
	}()

	wg.Wait() // goroutine이 완료될 때까지 대기

	assert.True(t, completed.Load(), "goroutine이 완료되었어야 함")
}

// TestChannelSolution - channel로 goroutine 완료 대기
func TestChannelSolution(t *testing.T) {
	done := make(chan struct{})

	go func() {
		defer close(done)
		time.Sleep(50 * time.Millisecond)
		// 작업 수행
	}()

	<-done // goroutine이 완료될 때까지 blocking

	t.Log("goroutine 완료 후 계속 진행")
}

// TestMultipleGoroutineWait - 여러 goroutine을 WaitGroup으로 대기
func TestMultipleGoroutineWait(t *testing.T) {
	const numWorkers = 5
	var wg sync.WaitGroup
	results := make([]int, numWorkers)

	wg.Add(numWorkers)
	for i := range numWorkers {
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			results[i] = i * 2
		}()
	}

	wg.Wait()

	for i, v := range results {
		assert.Equal(t, i*2, v)
	}
}
