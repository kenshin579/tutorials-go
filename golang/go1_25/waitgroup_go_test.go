package go1_25

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WaitGroup_기존방식_Add_Done(t *testing.T) {
	var wg sync.WaitGroup
	var counter atomic.Int64

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(5), counter.Load())
}

func Test_WaitGroup_새방식_Go(t *testing.T) {
	// Go 1.25: wg.Go()로 Add(1) + goroutine 생성 + Done() 자동화
	var wg sync.WaitGroup
	var counter atomic.Int64

	for i := 0; i < 5; i++ {
		wg.Go(func() {
			counter.Add(1)
		})
	}

	wg.Wait()
	assert.Equal(t, int64(5), counter.Load())
}

func Test_WaitGroup_Go_결과수집(t *testing.T) {
	var wg sync.WaitGroup
	results := make(chan int, 10)

	for i := range 10 {
		wg.Go(func() {
			results <- i * 2
		})
	}

	wg.Wait()
	close(results)

	var sum int
	for v := range results {
		sum += v
	}
	// 0+2+4+6+8+10+12+14+16+18 = 90
	assert.Equal(t, 90, sum)
}
