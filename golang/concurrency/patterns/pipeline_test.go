package patterns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Pipeline 패턴: 각 스테이지가 channel로 연결

// generator - 값을 생성하는 첫 번째 스테이지
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// square - 값을 제곱하는 중간 스테이지
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// filter - 조건을 만족하는 값만 통과시키는 스테이지
func filter(in <-chan int, predicate func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if predicate(n) {
				out <- n
			}
		}
	}()
	return out
}

// collect - channel의 모든 값을 slice로 수집
func collect(in <-chan int) []int {
	var result []int
	for v := range in {
		result = append(result, v)
	}
	return result
}

// TestPipeline - 파이프라인: generator → square → collect
func TestPipeline(t *testing.T) {
	// generator(1,2,3,4,5) → square → collect
	result := collect(square(generator(1, 2, 3, 4, 5)))

	assert.Equal(t, []int{1, 4, 9, 16, 25}, result)
}

// TestPipelineWithFilter - 파이프라인: generator → square → filter → collect
func TestPipelineWithFilter(t *testing.T) {
	// 제곱한 후 10 이상인 것만 필터링
	result := collect(
		filter(
			square(generator(1, 2, 3, 4, 5)),
			func(n int) bool { return n >= 10 },
		),
	)

	assert.Equal(t, []int{16, 25}, result)
}
