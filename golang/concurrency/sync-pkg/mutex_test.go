package sync_pkg

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRaceCondition - race condition 예시 (보호 없이 공유 변수 접근)
func TestRaceCondition(t *testing.T) {
	// 주의: 이 테스트는 -race 플래그 없이 실행해야 함 (race detector가 감지)
	// 여기서는 Mutex로 보호된 올바른 버전만 테스트
	counter := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	const numGoroutines = 1000
	wg.Add(numGoroutines)

	for range numGoroutines {
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	assert.Equal(t, numGoroutines, counter)
}

// TestMutexCriticalSection - Mutex로 임계 영역 보호
func TestMutexCriticalSection(t *testing.T) {
	type SafeCounter struct {
		mu sync.Mutex
		v  map[string]int
	}

	c := SafeCounter{v: make(map[string]int)}
	var wg sync.WaitGroup

	// 여러 goroutine에서 동시에 map에 접근
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.mu.Lock()
			c.v["key"]++
			c.mu.Unlock()
		}()
	}

	wg.Wait()

	c.mu.Lock()
	assert.Equal(t, 100, c.v["key"])
	c.mu.Unlock()
}

// TestRWMutex - RWMutex로 읽기 동시성 허용
func TestRWMutex(t *testing.T) {
	var rwmu sync.RWMutex
	data := map[string]string{"key": "value"}
	var wg sync.WaitGroup

	// 여러 reader 동시 실행 (RLock은 여러 goroutine이 동시에 획득 가능)
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rwmu.RLock()
			_ = data["key"]
			rwmu.RUnlock()
		}()
	}

	// writer는 exclusive lock (모든 reader/writer가 끝나야 획득)
	wg.Add(1)
	go func() {
		defer wg.Done()
		rwmu.Lock()
		data["key"] = "updated"
		rwmu.Unlock()
	}()

	wg.Wait()
}

// BenchmarkMutex - Mutex 벤치마크 (read 90% + write 10%)
func BenchmarkMutex(b *testing.B) {
	var mu sync.Mutex
	counter := 0

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 { // 10%는 write
				mu.Lock()
				counter++
				mu.Unlock()
			} else { // 90%는 read
				mu.Lock()
				_ = counter
				mu.Unlock()
			}
			i++
		}
	})
}

// BenchmarkRWMutex - RWMutex 벤치마크 (read 90% + write 10%)
func BenchmarkRWMutex(b *testing.B) {
	var rwmu sync.RWMutex
	counter := 0

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 { // 10%는 write
				rwmu.Lock()
				counter++
				rwmu.Unlock()
			} else { // 90%는 read
				rwmu.RLock()
				_ = counter
				rwmu.RUnlock()
			}
			i++
		}
	})
}
