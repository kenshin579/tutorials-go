package errhandling

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

// TestErrgroupBasic - errgroup 기본 사용법
func TestErrgroupBasic(t *testing.T) {
	g := new(errgroup.Group)

	g.Go(func() error {
		return nil // 성공
	})

	g.Go(func() error {
		return fmt.Errorf("task failed")
	})

	err := g.Wait() // 첫 번째 에러 반환
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task failed")
}

// TestErrgroupAllSuccess - 모든 작업 성공
func TestErrgroupAllSuccess(t *testing.T) {
	g := new(errgroup.Group)

	results := make([]int, 5)
	for i := range 5 {
		g.Go(func() error {
			results[i] = i * i
			return nil
		})
	}

	err := g.Wait()
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 4, 9, 16}, results)
}

// TestErrgroupWithContext - errgroup + context: 첫 에러 시 전체 취소
func TestErrgroupWithContext(t *testing.T) {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return fmt.Errorf("first error")
	})

	g.Go(func() error {
		<-ctx.Done() // 첫 번째 에러로 context가 취소됨
		return ctx.Err()
	})

	err := g.Wait()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "first error")
}

// TestErrgroupSetLimit - SetLimit로 동시성 제한
func TestErrgroupSetLimit(t *testing.T) {
	g := new(errgroup.Group)
	g.SetLimit(3) // 최대 3개 goroutine 동시 실행

	results := make([]int, 10)
	for i := range 10 {
		g.Go(func() error {
			results[i] = i
			return nil
		})
	}

	err := g.Wait()
	assert.NoError(t, err)
}
