package memory_model

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BenchmarkAtomicCounter - atomic 카운터 성능
func BenchmarkAtomicCounter(b *testing.B) {
	var counter atomic.Int64
	// RunParallel은 GOMAXPROCS만큼 goroutine을 띄워 실제 contention 상황을 측정한다
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// CPU instruction 한 번으로 처리되므로 lock 획득 비용이 없다
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
			// Lock/Unlock은 contention 시 goroutine parking까지 발생할 수 있어 atomic 대비 비싸다
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
			// 읽기 전용은 사실상 일반 메모리 read 수준 비용
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
			// 읽기만 해도 reader 사이에서 직렬화되므로 동시성이 떨어진다
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
			// RLock은 reader 간 병렬 허용이지만 내부 counter 갱신 비용은 남는다
			mu.RLock()
			_ = val
			mu.RUnlock()
		}
	})
}
