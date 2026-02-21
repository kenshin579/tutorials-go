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
	jobs := make(chan int, 10) // 작업을 분배할 공유 channel
	numWorkers := 3

	// worker별 결과 channel 생성
	workerResults := make([]chan int, numWorkers)
	for i := range numWorkers {
		workerResults[i] = make(chan int, 10)
	}

	// Fan-out: 각 worker가 같은 jobs channel에서 작업을 가져감
	var wg sync.WaitGroup
	for i := range numWorkers {
		wg.Add(1)
		go func() { // 각 worker goroutine
			defer wg.Done()
			for job := range jobs { // jobs가 close되면 루프 종료
				workerResults[i] <- job * job // 제곱 연산 후 결과 전송
			}
			close(workerResults[i])
		}()
	}

	// 9개의 작업을 channel에 전송
	for i := 1; i <= 9; i++ {
		jobs <- i
	}
	close(jobs) // 모든 작업 전송 완료 → worker들이 루프 종료

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

// fanIn - 여러 channel의 값을 하나의 channel로 합치는 함수
func fanIn(channels ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	merged := make(chan string) // 모든 결과가 모이는 단일 channel

	// 각 source channel마다 goroutine이 값을 merged로 전달
	for _, ch := range channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range ch { // source channel이 close되면 루프 종료
				merged <- v
			}
		}()
	}

	go func() {
		wg.Wait()     // 모든 source가 완료될 때까지 대기
		close(merged) // 모든 source 완료 후 merged channel 닫기
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
	for v := range merged { // merged channel이 close되면 루프 종료
		results = append(results, v)
	}

	assert.Len(t, results, 9) // 각 source 3개씩 = 총 9개
	t.Logf("merged results: %v", results)
}

// TestFanOutFanIn - fan-out + fan-in 조합으로 병렬 처리 파이프라인 구성
func TestFanOutFanIn(t *testing.T) {
	// generator: 슬라이스를 channel로 변환하는 헬퍼
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

	// square: 입력값을 제곱하여 문자열로 출력하는 worker
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

	// Fan-out: 입력을 3개 worker에 라운드로빈으로 분배
	input := generator(1, 2, 3, 4, 5, 6, 7, 8, 9)

	worker1In := make(chan int, 9)
	worker2In := make(chan int, 9)
	worker3In := make(chan int, 9)

	go func() {
		i := 0
		workers := []chan int{worker1In, worker2In, worker3In}
		for n := range input {
			workers[i%3] <- n // 라운드로빈으로 작업 분배
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
