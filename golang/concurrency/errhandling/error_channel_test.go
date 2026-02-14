package errhandling

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Result - 결과와 에러를 함께 전달하는 struct
type Result struct {
	Value int
	Err   error
}

// TestErrorChannel - error channel 패턴
func TestErrorChannel(t *testing.T) {
	work := func(id int) <-chan Result {
		ch := make(chan Result, 1)
		go func() {
			defer close(ch)
			if id%2 == 0 {
				ch <- Result{Value: id * 10, Err: nil}
			} else {
				ch <- Result{Err: fmt.Errorf("worker %d failed", id)}
			}
		}()
		return ch
	}

	// 3개의 worker 실행
	results := make([]<-chan Result, 3)
	for i := range 3 {
		results[i] = work(i)
	}

	var successes, failures int
	for _, ch := range results {
		r := <-ch
		if r.Err != nil {
			failures++
		} else {
			successes++
		}
	}

	assert.Equal(t, 2, successes) // 0, 2
	assert.Equal(t, 1, failures)  // 1
}

// TestMultiError - 여러 에러를 errors.Join으로 수집
func TestMultiError(t *testing.T) {
	var errs []error

	for i := range 5 {
		if i%2 != 0 {
			errs = append(errs, fmt.Errorf("task %d failed", i))
		}
	}

	combined := errors.Join(errs...)

	assert.Error(t, combined)
	assert.Contains(t, combined.Error(), "task 1 failed")
	assert.Contains(t, combined.Error(), "task 3 failed")
}
