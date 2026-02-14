package sync_pkg

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOnceBasic - sync.Once는 함수를 한 번만 실행
func TestOnceBasic(t *testing.T) {
	var once sync.Once
	var count atomic.Int64

	for range 100 {
		once.Do(func() {
			count.Add(1)
		})
	}

	assert.Equal(t, int64(1), count.Load())
}

// TestOnceSingleton - sync.Once로 singleton 패턴 구현
func TestOnceSingleton(t *testing.T) {
	type Config struct {
		DBHost string
		DBPort int
	}

	var (
		instance *Config
		once     sync.Once
	)

	getConfig := func() *Config {
		once.Do(func() {
			instance = &Config{
				DBHost: "localhost",
				DBPort: 5432,
			}
		})
		return instance
	}

	// 여러 goroutine에서 동시에 호출해도 같은 인스턴스 반환
	var wg sync.WaitGroup
	results := make([]*Config, 10)

	for i := range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = getConfig()
		}()
	}

	wg.Wait()

	// 모든 결과가 같은 포인터를 가리키는지 확인
	for i := 1; i < len(results); i++ {
		assert.Same(t, results[0], results[i])
	}
}

// TestOnceConcurrent - 여러 goroutine에서 동시에 Once.Do 호출
func TestOnceConcurrent(t *testing.T) {
	var once sync.Once
	var initCount atomic.Int64
	var wg sync.WaitGroup

	const numGoroutines = 100
	wg.Add(numGoroutines)

	for range numGoroutines {
		go func() {
			defer wg.Done()
			once.Do(func() {
				initCount.Add(1)
			})
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(1), initCount.Load())
}
