package patterns

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWorkerPool - Worker Pool 패턴: 고정 수의 worker가 job을 처리
func TestWorkerPool(t *testing.T) {
	type Job struct {
		ID    int
		Input int
	}
	type Result struct {
		JobID  int
		Output int
	}

	const numWorkers = 3
	const numJobs = 10

	// buffered channel: producer가 worker 속도에 묶이지 않고 미리 작업을 쌓아둘 수 있다
	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup
	for w := range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// jobs가 비면 블로킹 대기, close되면 자연스럽게 루프 종료
			for job := range jobs {
				results <- Result{
					JobID:  job.ID,
					Output: job.Input * job.Input,
				}
				t.Logf("worker %d processed job %d", w, job.ID)
			}
		}()
	}

	for i := range numJobs {
		jobs <- Job{ID: i, Input: i + 1}
	}
	// 종료 신호: close하지 않으면 worker들이 빈 channel에서 영원히 대기 → goroutine leak
	close(jobs)

	// results close는 별도 goroutine에서: 메인에서 wg.Wait()을 직접 부르면
	// 결과 수집 루프가 멈춰있어 worker의 송신이 블로킹 → deadlock
	go func() {
		wg.Wait()
		close(results)
	}()

	var collected []Result
	for r := range results {
		collected = append(collected, r)
	}

	assert.Len(t, collected, numJobs)
}

// TestWorkerPoolWithFunc - 함수형 Worker Pool
func TestWorkerPoolWithFunc(t *testing.T) {
	// processor 함수를 주입받아 작업 로직과 동시성 관리를 분리 (관심사 분리)
	workerPool := func(numWorkers int, jobs <-chan int, processor func(int) string) <-chan string {
		// unbuffered results: 수신자가 받을 때까지 worker가 송신에서 대기 (자연스러운 backpressure)
		results := make(chan string)
		var wg sync.WaitGroup

		for range numWorkers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// close(jobs) 시 range 종료 → worker도 종료
				for job := range jobs {
					results <- processor(job)
				}
			}()
		}

		// 모든 worker 종료 후 results close (수집 루프 종료를 위해 별도 goroutine에서)
		go func() {
			wg.Wait()
			close(results)
		}()

		return results // 수신 전용 반환: 호출자는 결과를 받기만 함
	}

	jobs := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		jobs <- i
	}
	close(jobs)

	// processor만 교체하면 동일한 worker pool로 다른 작업 처리 가능
	results := workerPool(3, jobs, func(n int) string {
		return fmt.Sprintf("%d^2=%d", n, n*n)
	})

	var collected []string
	for r := range results {
		collected = append(collected, r)
	}

	assert.Len(t, collected, 5)
}
