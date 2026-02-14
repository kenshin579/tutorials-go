package sync_pkg

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSyncMapBasic - sync.Map 기본 사용법
func TestSyncMapBasic(t *testing.T) {
	var m sync.Map

	// Store: 값 저장
	m.Store("key1", "value1")
	m.Store("key2", 42)

	// Load: 값 조회
	val, ok := m.Load("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	// LoadOrStore: 없으면 저장, 있으면 기존 값 반환
	actual, loaded := m.LoadOrStore("key3", "new")
	assert.False(t, loaded) // 새로 저장됨
	assert.Equal(t, "new", actual)

	actual, loaded = m.LoadOrStore("key3", "another")
	assert.True(t, loaded) // 이미 존재
	assert.Equal(t, "new", actual)

	// Delete: 삭제
	m.Delete("key2")
	_, ok = m.Load("key2")
	assert.False(t, ok)
}

// TestSyncMapRange - sync.Map 순회
func TestSyncMapRange(t *testing.T) {
	var m sync.Map

	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	count := 0
	m.Range(func(key, value any) bool {
		count++
		return true // true: 계속 순회, false: 중단
	})

	assert.Equal(t, 3, count)
}

// TestSyncMapConcurrent - sync.Map은 concurrent-safe
func TestSyncMapConcurrent(t *testing.T) {
	var m sync.Map
	var wg sync.WaitGroup

	// 여러 goroutine에서 동시에 읽고 쓰기
	const numGoroutines = 100
	wg.Add(numGoroutines * 2)

	for i := range numGoroutines {
		// writer
		go func() {
			defer wg.Done()
			m.Store(fmt.Sprintf("key-%d", i), i)
		}()

		// reader
		go func() {
			defer wg.Done()
			m.Load(fmt.Sprintf("key-%d", i))
		}()
	}

	wg.Wait()
}

// BenchmarkMapWithMutex - 일반 map + Mutex
func BenchmarkMapWithMutex(b *testing.B) {
	var mu sync.RWMutex
	m := make(map[int]int)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				mu.Lock()
				m[i] = i
				mu.Unlock()
			} else {
				mu.RLock()
				_ = m[i]
				mu.RUnlock()
			}
			i++
		}
	})
}

// BenchmarkSyncMap - sync.Map
func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				m.Store(i, i)
			} else {
				m.Load(i)
			}
			i++
		}
	})
}
