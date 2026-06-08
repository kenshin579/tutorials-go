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

	// g.Go: 내부적으로 wg.Add(1) + goroutine 실행 + 에러 저장까지 자동 처리
	g.Go(func() error {
		return nil // 성공
	})

	g.Go(func() error {
		return fmt.Errorf("task failed")
	})

	err := g.Wait() // 모든 goroutine 완료 대기 + 첫 번째 non-nil 에러 반환
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task failed")
}

// TestErrgroupAllSuccess - 모든 작업 성공
func TestErrgroupAllSuccess(t *testing.T) {
	g := new(errgroup.Group)

	// 각 goroutine이 서로 다른 인덱스에만 쓰므로 mutex 없이 안전 (race 없음)
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
	// WithContext: 첫 에러 반환 시 ctx 자동 cancel → 다른 goroutine에 취소 신호 전파
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return fmt.Errorf("first error")
	})

	g.Go(func() error {
		<-ctx.Done() // 위 goroutine의 에러로 ctx가 취소되어 깨어남
		return ctx.Err()
	})

	// Wait는 g.Go에 등록된 함수가 반환한 첫 번째 non-nil 에러를 돌려준다
	err := g.Wait()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "first error")
}

// TestErrgroupSetLimit - SetLimit로 동시성 제한
func TestErrgroupSetLimit(t *testing.T) {
	g := new(errgroup.Group)
	// 내부적으로 buffered channel 기반 세마포어 → 슬롯이 빌 때까지 g.Go가 블로킹
	g.SetLimit(3)

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
