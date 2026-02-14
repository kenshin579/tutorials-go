package memory_model

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BenchmarkAtomicCounter - atomic 카운터 성능
func BenchmarkAtomicCounter(b *testing.B) {
	var counter atomic.Int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Add(1)
		}
	})
}

// BenchmarkMutexCounter - mutex 카운터 성능
func BenchmarkMutexCounter(b *testing.B) {
	var mu sync.Mutex
	var counter int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

// BenchmarkAtomicLoad - atomic Load 성능
func BenchmarkAtomicLoad(b *testing.B) {
	var val atomic.Int64
	val.Store(42)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = val.Load()
		}
	})
}

// BenchmarkMutexRead - mutex 읽기 성능
func BenchmarkMutexRead(b *testing.B) {
	var mu sync.Mutex
	val := int64(42)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			_ = val
			mu.Unlock()
		}
	})
}

// BenchmarkRWMutexRead - RWMutex 읽기 성능
func BenchmarkRWMutexRead(b *testing.B) {
	var mu sync.RWMutex
	val := int64(42)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.RLock()
			_ = val
			mu.RUnlock()
		}
	})
}
