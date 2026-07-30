package bloomfilter

import (
	"fmt"
	"runtime"
	"testing"
)

// makeURL은 약 54바이트짜리 URL을 만든다.
func makeURL(i int) string {
	return fmt.Sprintf("https://example.com/page/%09d/detail-view-section", i)
}

// makeKey는 length 바이트짜리 고유한 키를 만든다. length는 10 이상이어야 한다.
func makeKey(i, length int) string {
	s := fmt.Sprintf("%09d-", i) // 앞 10바이트가 항상 고유하다
	for len(s) < length {
		s += "x"
	}
	return s[:length]
}

// 본문 3.3절 표의 근거. 목표 FPR별로 필요한 비트 수와 메모리를 출력한다.
func TestTable_목표FPR별_메모리(t *testing.T) {
	const n = 1_000_000

	t.Log("목표p | m(비트) | 메모리 | k | 원소당비트")
	for _, p := range []float64{0.1, 0.01, 0.001, 0.0001} {
		m := OptimalM(n, p)
		k := OptimalK(m, n)
		sizeBytes := (m + 63) / 64 * 8 // 실제 할당되는 워드 배열 크기
		t.Logf("%g | %d | %.2f MiB | %d | %.2f",
			p, m, float64(sizeBytes)/(1024*1024), k, float64(m)/float64(n))
	}
}

// 본문 3.4절 표의 근거.
// map은 키를 계속 들고 있어야 하지만 Bloom Filter는 원소를 저장하지 않는다.
// 그래서 배수는 원소가 길수록 커진다 - "99배"는 특정 키 크기에서만 나오는 값이다.
func TestMemory_map과_BloomFilter_비교(t *testing.T) {
	const n = 1_000_000

	// Bloom Filter 쪽은 원소 크기와 무관하게 일정하다
	bloomBytes := measureHeap(func() any {
		f := NewWithEstimates(n, 0.01)
		for i := 0; i < n; i++ {
			f.Add([]byte(makeURL(i)))
		}
		return f
	})
	t.Logf("원소 수: %d", n)
	t.Logf("Bloom Filter(p=0.01): %.2f MiB (원소 크기와 무관)", float64(bloomBytes)/(1024*1024))

	intBytes := measureHeap(func() any {
		m := make(map[uint64]struct{})
		for i := 0; i < n; i++ {
			m[uint64(i)] = struct{}{}
		}
		return m
	})
	t.Logf("map[uint64]struct{} (8바이트): %.2f MiB, 배수=%.1f배",
		float64(intBytes)/(1024*1024), float64(intBytes)/float64(bloomBytes))

	for _, length := range []int{15, 54, 200} {
		mapBytes := measureHeap(func() any {
			m := make(map[string]struct{})
			for i := 0; i < n; i++ {
				m[makeKey(i, length)] = struct{}{}
			}
			return m
		})
		t.Logf("map[string]struct{} (%d바이트): %.2f MiB, 배수=%.1f배",
			length, float64(mapBytes)/(1024*1024), float64(mapBytes)/float64(bloomBytes))
	}
}

// measureHeap는 build가 만든 자료구조가 차지하는 힙 증가량을 잰다.
// 반환값을 KeepAlive로 살려두어야 GC가 측정 전에 회수하지 않는다.
func measureHeap(build func() any) uint64 {
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)

	held := build()

	runtime.GC()
	runtime.ReadMemStats(&after)
	runtime.KeepAlive(held)

	return after.HeapAlloc - before.HeapAlloc
}
