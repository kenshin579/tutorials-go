package debugging

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRaceConditionFixed - race condition을 mutex로 수정한 버전
// 원래 코드: counter++를 여러 goroutine에서 동시 실행 → race condition
// 수정: sync.Mutex로 임계영역 보호
func TestRaceConditionFixed(t *testing.T) {
	var mu sync.Mutex
	counter := 0
	var wg sync.WaitGroup

	wg.Add(1000)
	for range 1000 {
		go func() {
			defer wg.Done()
			// critical section은 최소화해야 경합으로 인한 대기 시간을 줄일 수 있다
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	assert.Equal(t, 1000, counter)
}

// TestRaceConditionAtomicFix - atomic으로 race condition 수정
func TestRaceConditionAtomicFix(t *testing.T) {
	var counter atomic.Int64
	var wg sync.WaitGroup

	wg.Add(1000)
	for range 1000 {
		go func() {
			defer wg.Done()
			// CPU atomic instruction을 사용 → lock 경합 없이 단순 증가에 최적
			counter.Add(1)
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(1000), counter.Load())
}

// TestMapRaceFixed - map의 concurrent access를 sync.Map으로 수정
// 원래 코드: 일반 map을 여러 goroutine에서 동시 읽기/쓰기 → fatal: concurrent map writes
// 수정: sync.Map 사용
func TestMapRaceFixed(t *testing.T) {
	// sync.Map: 읽기 빈번 + 쓰기 드문 워크로드에 최적
	// 일반 map은 concurrent write 시 fatal error로 즉시 종료됨
	var m sync.Map
	var wg sync.WaitGroup

	// 동시에 쓰기
	wg.Add(100)
	for i := range 100 {
		go func() {
			defer wg.Done()
			m.Store(i, i*10)
		}()
	}

	// 동시에 읽기
	wg.Add(100)
	for i := range 100 {
		go func() {
			defer wg.Done()
			m.Load(i)
		}()
	}

	wg.Wait()
	t.Log("sync.Map으로 안전한 concurrent access 완료")
}

// TestSliceRaceFixed - slice의 concurrent access 수정
// 원래 코드: slice append를 여러 goroutine에서 동시 실행 → race condition
// 수정: 인덱스별 독립 접근 (각 goroutine이 고유 인덱스에만 씀)
func TestSliceRaceFixed(t *testing.T) {
	// 미리 길이를 할당한 slice → 각 요소는 독립 메모리 주소
	results := make([]int, 100)
	var wg sync.WaitGroup

	wg.Add(100)
	for i := range 100 {
		go func() {
			defer wg.Done()
			// 각 goroutine이 서로 다른 인덱스에만 쓰므로 race가 아님
			// 단, append처럼 slice 구조 자체를 수정하는 연산은 race 발생
			results[i] = i * 2
		}()
	}

	wg.Wait()
	assert.Equal(t, 0, results[0])
	assert.Equal(t, 198, results[99])
}
