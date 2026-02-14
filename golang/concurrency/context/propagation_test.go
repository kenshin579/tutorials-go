package context_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestPropagation - context 전파 체인: parent 취소 → 모든 child 취소
func TestPropagation(t *testing.T) {
	root, rootCancel := context.WithCancel(context.Background())

	var stopped atomic.Int64

	// 3개의 child goroutine 생성
	for i := range 3 {
		child, childCancel := context.WithCancel(root)
		defer childCancel()

		go func() {
			<-child.Done()
			stopped.Add(1)
			t.Logf("child %d stopped", i)
		}()
	}

	rootCancel() // root 취소 → 모든 child 취소
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, int64(3), stopped.Load())
}

// worker - context를 첫 번째 파라미터로 받는 함수 (Go 관례)
func worker(ctx context.Context, id int, results chan<- string) {
	for {
		select {
		case <-ctx.Done():
			results <- fmt.Sprintf("worker %d: stopped", id)
			return
		case <-time.After(10 * time.Millisecond):
			results <- fmt.Sprintf("worker %d: working", id)
		}
	}
}

// TestContextAsFirstParam - context는 함수의 첫 번째 파라미터로 전달
func TestContextAsFirstParam(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	results := make(chan string, 100)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); worker(ctx, 1, results) }()
	go func() { defer wg.Done(); worker(ctx, 2, results) }()

	// worker가 모두 종료된 후 channel close
	go func() {
		wg.Wait()
		close(results)
	}()

	var messages []string
	for msg := range results {
		messages = append(messages, msg)
	}

	t.Logf("messages: %v", messages)
	assert.Greater(t, len(messages), 0)
}
