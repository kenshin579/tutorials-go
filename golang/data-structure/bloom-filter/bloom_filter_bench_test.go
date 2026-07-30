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

// missKeys는 필터에 넣지 않은 키를 만든다. 채운 범위(0..benchN-1) 밖을 쓴다.
func missKeys(count int) [][]byte {
	keys := make([][]byte, count)
	for i := range keys {
		keys[i] = []byte(makeURL(benchN + i))
	}
	return keys
}

// 조회 벤치마크는 설계 용량까지 채운 필터로 측정한다.
// 필터가 비어 있으면 "없는 키"가 첫 프로브에서 곧바로 반환되어
// 실제보다 빠르게 나오고 구현 간 차이도 사라진다.
func fillMine() *BloomFilter {
	f := NewWithEstimates(benchN, 0.01)
	for i := 0; i < benchN; i++ {
		f.Add([]byte(makeURL(i)))
	}
	return f
}

func fillLib() *bloom.BloomFilter {
	f := bloom.NewWithEstimates(benchN, 0.01)
	for i := 0; i < benchN; i++ {
		f.Add([]byte(makeURL(i)))
	}
	return f
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
	f := fillMine()
	keys := benchKeys(1000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Contains(keys[i%len(keys)])
	}
}

func BenchmarkContains_라이브러리(b *testing.B) {
	f := fillLib()
	keys := benchKeys(1000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Test(keys[i%len(keys)])
	}
}

func BenchmarkContains_직접구현_없는키(b *testing.B) {
	f := fillMine()
	keys := missKeys(1000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Contains(keys[i%len(keys)])
	}
}

func BenchmarkContains_라이브러리_없는키(b *testing.B) {
	f := fillLib()
	keys := missKeys(1000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Test(keys[i%len(keys)])
	}
}
