package go_testing

import (
	"fmt"
	"testing"
)

// BenchmarkAverage - 기본 벤치마크
// 실행: go test -bench=BenchmarkAverage -benchmem
func BenchmarkAverage(b *testing.B) {
	nos := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for b.Loop() {
		Average(nos...)
	}
}

// BenchmarkAverage_ResetTimer - b.ResetTimer() 활용
// 셋업 시간을 벤치마크 측정에서 제외한다.
func BenchmarkAverage_ResetTimer(b *testing.B) {
	// 셋업: 큰 슬라이스 생성 (측정 대상 아님)
	nos := make([]int, 1000)
	for i := range nos {
		nos[i] = i + 1
	}
	b.ResetTimer() // 여기서부터 측정 시작

	for b.Loop() {
		Average(nos...)
	}
}

// BenchmarkAverage_ReportAllocs - b.ReportAllocs() 활용
// -benchmem 플래그 없이도 메모리 할당 정보를 출력한다.
func BenchmarkAverage_ReportAllocs(b *testing.B) {
	b.ReportAllocs()
	nos := []int{1, 2, 3, 4, 5}
	for b.Loop() {
		Average(nos...)
	}
}

// BenchmarkAverage_SubBenchmark - b.Run()을 사용한 서브벤치마크
// 입력 크기별로 성능을 비교할 수 있다.
// 실행: go test -bench=BenchmarkAverage_SubBenchmark -benchmem
// 특정 크기만: go test -bench=BenchmarkAverage_SubBenchmark/size=100
func BenchmarkAverage_SubBenchmark(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		nos := make([]int, size)
		for i := range nos {
			nos[i] = i + 1
		}

		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			for b.Loop() {
				Average(nos...)
			}
		})
	}
}

// generateSlice - 벤치마크 헬퍼 함수
func generateSlice(size int) []int {
	nos := make([]int, size)
	for i := range nos {
		nos[i] = i + 1
	}
	return nos
}

// BenchmarkAverage_MemoryComparison - 슬라이스 복사 vs 직접 전달 비교
// 서브벤치마크로 두 가지 접근법의 성능 차이를 측정한다.
func BenchmarkAverage_MemoryComparison(b *testing.B) {
	nos := generateSlice(1000)

	b.Run("direct", func(b *testing.B) {
		for b.Loop() {
			Average(nos...)
		}
	})

	b.Run("copy", func(b *testing.B) {
		for b.Loop() {
			copied := make([]int, len(nos))
			copy(copied, nos)
			Average(copied...)
		}
	})
}
