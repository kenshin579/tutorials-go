package sync_pkg

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWaitGroupBasic - WaitGroup 기본 사용법
func TestWaitGroupBasic(t *testing.T) {
	var wg sync.WaitGroup
	var results []int
	var mu sync.Mutex

	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			results = append(results, i)
			mu.Unlock()
		}()
	}

	wg.Wait() // 모든 goroutine 완료 대기
	assert.Len(t, results, 5)
}

// TestWaitGroupCounter - WaitGroup으로 병렬 작업 후 결과 집계
func TestWaitGroupCounter(t *testing.T) {
	var wg sync.WaitGroup
	var counter atomic.Int64

	const numWorkers = 100
	wg.Add(numWorkers)

	for range numWorkers {
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(numWorkers), counter.Load())
}

// TestWaitGroupAddBeforeGo - wg.Add는 반드시 go 문 전에 호출해야 한다
func TestWaitGroupAddBeforeGo(t *testing.T) {
	var wg sync.WaitGroup
	var counter atomic.Int64

	for range 10 {
		wg.Add(1) // go 문 전에 Add 호출 (올바른 패턴)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(10), counter.Load())
}
