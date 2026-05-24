package memory_model

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAtomicInt64 - atomic.Int64 기본 연산
func TestAtomicInt64(t *testing.T) {
	// atomic.Int64는 64비트 정수를 lock 없이 원자적으로 처리
	// (32비트 플랫폼에서 발생할 수 있는 word tearing 문제도 회피)
	var counter atomic.Int64

	counter.Store(10)
	assert.Equal(t, int64(10), counter.Load())

	counter.Add(5) // 단일 CPU instruction (LOCK XADD)으로 처리되어 race-free
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

	// Swap은 "값 변경 + 이전 값 확인"을 한 번에 처리해서
	// "최초로 true를 설정한 자가 누구인가" 같은 단일-실행 패턴에 유용
	old := flag.Swap(false)
	assert.True(t, old)
	assert.False(t, flag.Load())
}

// TestAtomicCompareAndSwap - CAS (Compare-And-Swap) 연산
func TestAtomicCompareAndSwap(t *testing.T) {
	var counter atomic.Int64
	counter.Store(100)

	// CAS는 "비교 + 교체"를 원자적으로 수행 → race 없이 조건부 업데이트 가능
	swapped := counter.CompareAndSwap(100, 200)
	assert.True(t, swapped)
	assert.Equal(t, int64(200), counter.Load())

	// 실패는 "다른 누군가가 먼저 값을 바꿨다"는 신호이므로 호출자가 재시도를 결정한다
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
			// lock-free 알고리즘의 표준 패턴: 읽기 → 판단 → CAS → 실패 시 재시도
			for {
				current := max.Load()
				if v <= current {
					break // 이미 더 큰 값이 있으면 종료 (불필요한 CAS 회피)
				}
				if max.CompareAndSwap(current, v) {
					break
				}
				// CAS 실패 = current를 읽은 후 누가 먼저 바꿨다는 뜻 → 새 값으로 다시 시도
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

	// atomic.Value는 struct 전체를 한 번에 교체하므로 부분 갱신된 중간 상태가 노출되지 않는다
	var config atomic.Value

	config.Store(Config{MaxConns: 10, Timeout: 30})

	loaded := config.Load().(Config)
	assert.Equal(t, 10, loaded.MaxConns)
	assert.Equal(t, 30, loaded.Timeout)

	// 필드 단위 수정이 아니라 "새 Config 통째로 교체" 방식 → reader는 항상 일관된 스냅샷을 본다
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

	// writer: hot-reload처럼 config를 계속 교체하는 쪽
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 2; i <= 100; i++ {
			config.Store(Config{Version: i})
		}
	}()

	// reader: writer와 동시에 동작해도 partial write가 보이지 않는다 (Version > 0 항상 보장)
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
			// 1000개 goroutine이 동시에 호출해도 Add는 atomic이라 update 손실 없음
			// (Mutex 없이 read-modify-write가 안전하게 처리되는 이유)
			counter.Add(1)
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(1000), counter.Load())
}
