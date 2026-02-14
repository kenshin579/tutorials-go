package _select

import (
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFanOut - 하나의 입력을 여러 worker에게 분배
func TestFanOut(t *testing.T) {
	jobs := make(chan int, 10)
	numWorkers := 3

	// 결과를 모을 channels (worker별 하나)
	workerResults := make([]chan int, numWorkers)
	for i := range numWorkers {
		workerResults[i] = make(chan int, 10)
	}

	// Fan-out: 각 worker가 jobs channel에서 작업을 가져감
	var wg sync.WaitGroup
	for i := range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				workerResults[i] <- job * job // 제곱 계산
			}
			close(workerResults[i])
		}()
	}

	// 작업 투입
	for i := 1; i <= 9; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	// 결과 수집
	var results []int
	for _, ch := range workerResults {
		for v := range ch {
			results = append(results, v)
		}
	}

	sort.Ints(results)
	assert.Equal(t, []int{1, 4, 9, 16, 25, 36, 49, 64, 81}, results)
}

// fanIn - 여러 channel을 하나로 합치는 함수
func fanIn(channels ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	merged := make(chan string)

	// 각 channel에서 값을 읽어 merged channel로 전달
	for _, ch := range channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range ch {
				merged <- v
			}
		}()
	}

	// 모든 channel이 닫히면 merged도 닫기
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// TestFanIn - 여러 channel의 결과를 하나로 합치기
func TestFanIn(t *testing.T) {
	// 3개의 독립적인 데이터 소스
	source1 := make(chan string, 3)
	source2 := make(chan string, 3)
	source3 := make(chan string, 3)

	go func() {
		for _, s := range []string{"a1", "a2", "a3"} {
			source1 <- s
		}
		close(source1)
	}()

	go func() {
		for _, s := range []string{"b1", "b2", "b3"} {
			source2 <- s
		}
		close(source2)
	}()

	go func() {
		for _, s := range []string{"c1", "c2", "c3"} {
			source3 <- s
		}
		close(source3)
	}()

	// Fan-in: 3개 channel을 하나로 합침
	merged := fanIn(source1, source2, source3)

	var results []string
	for v := range merged {
		results = append(results, v)
	}

	assert.Len(t, results, 9)
	t.Logf("merged results: %v", results)
}

// TestFanOutFanIn - fan-out + fan-in 조합
func TestFanOutFanIn(t *testing.T) {
	// 작업 생성 (generator)
	generator := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}

	// worker: 입력을 제곱하여 출력
	square := func(in <-chan int) <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			for n := range in {
				time.Sleep(5 * time.Millisecond) // 작업 시뮬레이션
				out <- fmt.Sprintf("%d^2=%d", n, n*n)
			}
		}()
		return out
	}

	// Fan-out: 입력을 3개 worker에 분배
	input := generator(1, 2, 3, 4, 5, 6, 7, 8, 9)

	// 입력을 버퍼링하여 worker에게 분배
	worker1In := make(chan int, 9)
	worker2In := make(chan int, 9)
	worker3In := make(chan int, 9)

	go func() {
		i := 0
		workers := []chan int{worker1In, worker2In, worker3In}
		for n := range input {
			workers[i%3] <- n
			i++
		}
		close(worker1In)
		close(worker2In)
		close(worker3In)
	}()

	w1 := square(worker1In)
	w2 := square(worker2In)
	w3 := square(worker3In)

	// Fan-in: 3개 worker 결과를 하나로 합침
	merged := fanIn(w1, w2, w3)

	var results []string
	for v := range merged {
		results = append(results, v)
	}

	assert.Len(t, results, 9)
	t.Logf("fan-out/fan-in results: %v", results)
}
