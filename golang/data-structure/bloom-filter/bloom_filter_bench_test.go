package bloomfilter

import (
	"fmt"
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
)

const benchN = 1_000_000

func benchKeys(count int) [][]byte {
	keys := make([][]byte, count)
	for i := range keys {
		keys[i] = []byte(fmt.Sprintf("key-%d", i))
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
