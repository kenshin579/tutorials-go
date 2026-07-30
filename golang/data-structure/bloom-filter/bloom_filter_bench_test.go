package bloomfilter

import (
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
)

const benchN = 1_000_000

func benchKeys(count int) [][]byte {
	keys := make([][]byte, count)
	for i := range keys {
		keys[i] = []byte(makeURL(i)) // 글 전체가 URL 중복 제거 시나리오이므로 통일한다
	}
	return keys
}

func BenchmarkAdd_직접구현(b *testing.B) {
	keys := benchKeys(1000)
	f := NewWithEstimates(benchN, 0.01)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f.Add(keys[i%len(keys)])
	}
}

func BenchmarkAdd_라이브러리(b *testing.B) {
	keys := benchKeys(1000)
	f := bloom.NewWithEstimates(benchN, 0.01)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f.Add(keys[i%len(keys)])
	}
}

func BenchmarkContains_직접구현(b *testing.B) {
	keys := benchKeys(1000)
	f := NewWithEstimates(benchN, 0.01)
	for _, k := range keys {
		f.Add(k)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f.Contains(keys[i%len(keys)])
	}
}

func BenchmarkContains_라이브러리(b *testing.B) {
	keys := benchKeys(1000)
	f := bloom.NewWithEstimates(benchN, 0.01)
	for _, k := range keys {
		f.Add(k)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f.Test(keys[i%len(keys)])
	}
}

// 실제 워크로드는 대부분 "없는 키" 조회다. 첫 0비트에서 조기 반환하므로 더 빠르다.
func BenchmarkContains_직접구현_없는키(b *testing.B) {
	keys := benchKeys(1000)
	missing := []byte("이 필터에-절대-없는-키")
	f := NewWithEstimates(benchN, 0.01)
	for _, k := range keys {
		f.Add(k)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Contains(missing)
	}
}

func BenchmarkContains_라이브러리_없는키(b *testing.B) {
	keys := benchKeys(1000)
	missing := []byte("이 필터에-절대-없는-키")
	f := bloom.NewWithEstimates(benchN, 0.01)
	for _, k := range keys {
		f.Add(k)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Test(missing)
	}
}
