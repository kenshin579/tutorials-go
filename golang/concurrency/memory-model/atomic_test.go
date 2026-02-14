package memory_model

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAtomicInt64 - atomic.Int64 기본 연산
func TestAtomicInt64(t *testing.T) {
	var counter atomic.Int64

	counter.Store(10)
	assert.Equal(t, int64(10), counter.Load())

	counter.Add(5)
	assert.Equal(t, int64(15), counter.Load())

	counter.Add(-3)
	assert.Equal(t, int64(12), counter.Load())
}

// TestAtomicBool - atomic.Bool 플래그
func TestAtomicBool(t *testing.T) {
	var flag atomic.Bool

	assert.False(t, flag.Load())

	flag.Store(true)
	assert.True(t, flag.Load())

	// Swap: 이전 값 반환
	old := flag.Swap(false)
	assert.True(t, old)
	assert.False(t, flag.Load())
}

// TestAtomicCompareAndSwap - CAS (Compare-And-Swap) 연산
func TestAtomicCompareAndSwap(t *testing.T) {
	var counter atomic.Int64
	counter.Store(100)

	// 현재 값이 100이면 200으로 교체 → 성공
	swapped := counter.CompareAndSwap(100, 200)
	assert.True(t, swapped)
	assert.Equal(t, int64(200), counter.Load())

	// 현재 값이 100이 아니므로 교체 실패
	swapped = counter.CompareAndSwap(100, 300)
	assert.False(t, swapped)
	assert.Equal(t, int64(200), counter.Load()) // 여전히 200
}

// TestAtomicCASLoop - CAS 루프 패턴 (lock-free 업데이트)
func TestAtomicCASLoop(t *testing.T) {
	var max atomic.Int64

	var wg sync.WaitGroup
	values := []int64{5, 3, 8, 1, 9, 2, 7}

	wg.Add(len(values))
	for _, v := range values {
		go func() {
			defer wg.Done()
			for {
				current := max.Load()
				if v <= current {
					break // 현재 값이 이미 크면 업데이트 불필요
				}
				if max.CompareAndSwap(current, v) {
					break // 성공적으로 업데이트
				}
				// CAS 실패 시 다시 시도 (다른 goroutine이 먼저 변경)
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(9), max.Load())
}

// TestAtomicValue - atomic.Value로 임의 타입 저장
func TestAtomicValue(t *testing.T) {
	type Config struct {
		MaxConns int
		Timeout  int
	}

	var config atomic.Value

	// 초기 설정 저장
	config.Store(Config{MaxConns: 10, Timeout: 30})

	loaded := config.Load().(Config)
	assert.Equal(t, 10, loaded.MaxConns)
	assert.Equal(t, 30, loaded.Timeout)

	// 설정 업데이트 (전체 교체)
	config.Store(Config{MaxConns: 20, Timeout: 60})

	updated := config.Load().(Config)
	assert.Equal(t, 20, updated.MaxConns)
	assert.Equal(t, 60, updated.Timeout)
}

// TestAtomicValueConcurrent - atomic.Value를 사용한 concurrent config 갱신
func TestAtomicValueConcurrent(t *testing.T) {
	type Config struct {
		Version int
	}

	var config atomic.Value
	config.Store(Config{Version: 1})

	var wg sync.WaitGroup

	// writer: 설정 업데이트
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 2; i <= 100; i++ {
			config.Store(Config{Version: i})
		}
	}()

	// reader: 설정 읽기
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 100 {
			cfg := config.Load().(Config)
			assert.Greater(t, cfg.Version, 0)
		}
	}()

	wg.Wait()

	final := config.Load().(Config)
	assert.Equal(t, 100, final.Version)
}

// TestAtomicCounter - 여러 goroutine에서 atomic counter 증가
func TestAtomicCounter(t *testing.T) {
	var counter atomic.Int64
	var wg sync.WaitGroup

	wg.Add(1000)
	for range 1000 {
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(1000), counter.Load())
}
