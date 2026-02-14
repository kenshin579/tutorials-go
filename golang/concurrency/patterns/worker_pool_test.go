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

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	// Worker 시작
	var wg sync.WaitGroup
	for w := range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				// 작업 처리: 제곱 계산
				results <- Result{
					JobID:  job.ID,
					Output: job.Input * job.Input,
				}
				t.Logf("worker %d processed job %d", w, job.ID)
			}
		}()
	}

	// Job 투입
	for i := range numJobs {
		jobs <- Job{ID: i, Input: i + 1}
	}
	close(jobs)

	// Worker 완료 후 results channel close
	go func() {
		wg.Wait()
		close(results)
	}()

	// 결과 수집
	var collected []Result
	for r := range results {
		collected = append(collected, r)
	}

	assert.Len(t, collected, numJobs)
}

// TestWorkerPoolWithFunc - 함수형 Worker Pool
func TestWorkerPoolWithFunc(t *testing.T) {
	workerPool := func(numWorkers int, jobs <-chan int, processor func(int) string) <-chan string {
		results := make(chan string)
		var wg sync.WaitGroup

		for range numWorkers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for job := range jobs {
					results <- processor(job)
				}
			}()
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		return results
	}

	jobs := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		jobs <- i
	}
	close(jobs)

	results := workerPool(3, jobs, func(n int) string {
		return fmt.Sprintf("%d^2=%d", n, n*n)
	})

	var collected []string
	for r := range results {
		collected = append(collected, r)
	}

	assert.Len(t, collected, 5)
}
