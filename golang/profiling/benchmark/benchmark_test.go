package benchmark

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

// 문자열 연결 방식 비교를 통한 벤치마크 프로파일링 예제
//
// 프로파일 수집 명령어:
//   go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof -benchmem
//
// 프로파일 분석:
//   go tool pprof cpu.prof
//   go tool pprof -http=:8080 mem.prof

func BenchmarkConcatPlus(b *testing.B) {
	for b.Loop() {
		concatPlus(1000)
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	for b.Loop() {
		concatBuilder(1000)
	}
}

func BenchmarkConcatSprintf(b *testing.B) {
	for b.Loop() {
		concatSprintf(1000)
	}
}

func BenchmarkHashComputation(b *testing.B) {
	data := []byte(strings.Repeat("hello world", 1000))
	b.ResetTimer()
	for b.Loop() {
		hashComputation(data)
	}
}

// concatPlus는 + 연산자로 문자열을 연결한다 (비효율적).
func concatPlus(n int) string {
	s := ""
	for i := range n {
		s += fmt.Sprintf("item-%d ", i)
	}
	return s
}

// concatBuilder는 strings.Builder로 문자열을 연결한다 (효율적).
func concatBuilder(n int) string {
	var b strings.Builder
	for i := range n {
		fmt.Fprintf(&b, "item-%d ", i)
	}
	return b.String()
}

// concatSprintf는 Sprintf로 문자열을 연결한다 (비효율적).
func concatSprintf(n int) string {
	s := ""
	for i := range n {
		s = fmt.Sprintf("%s item-%d", s, i)
	}
	return s
}

// hashComputation은 CPU 부하를 생성하는 해시 계산 함수이다.
func hashComputation(data []byte) [32]byte {
	result := sha256.Sum256(data)
	for range 100 {
		result = sha256.Sum256(result[:])
	}
	return result
}
