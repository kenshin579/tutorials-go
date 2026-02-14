package patterns

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSemaphore - buffered channel로 동시 실행 수 제한
func TestSemaphore(t *testing.T) {
	const maxConcurrency = 3
	sem := make(chan struct{}, maxConcurrency)

	var maxConcurrent atomic.Int64
	var currentConcurrent atomic.Int64
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem <- struct{}{}        // 세마포어 획득
			defer func() { <-sem }() // 세마포어 해제

			// 동시 실행 수 추적
			cur := currentConcurrent.Add(1)
			for {
				old := maxConcurrent.Load()
				if cur <= old || maxConcurrent.CompareAndSwap(old, cur) {
					break
				}
			}

			time.Sleep(20 * time.Millisecond) // 작업 시뮬레이션
			currentConcurrent.Add(-1)
		}()
	}

	wg.Wait()

	t.Logf("최대 동시 실행 수: %d", maxConcurrent.Load())
	assert.LessOrEqual(t, maxConcurrent.Load(), int64(maxConcurrency))
}
