// Package bloomfilter는 Bloom Filter를 비트 배열과 이중 해싱으로 직접 구현한 예제이다.
package bloomfilter

import (
	"math"

	"github.com/cespare/xxhash/v2"
)

// BloomFilter는 비트 배열 기반의 확률적 집합이다.
// Contains가 false를 반환하면 원소는 확실히 없고,
// true를 반환하면 있을 수도 있다(false positive).
type BloomFilter struct {
	bits []uint64 // 비트 배열을 uint64 워드 단위로 저장
	m    uint64   // 전체 비트 수
	k    uint64   // 해시 함수 개수
	n    uint64   // 추가된 원소 수
}

// New는 비트 수 m과 해시 함수 개수 k를 직접 지정해 생성한다.
func New(m, k uint64) *BloomFilter {
	if m == 0 {
		m = 1
	}
	if k == 0 {
		k = 1
	}
	return &BloomFilter{
		bits: make([]uint64, (m+63)/64), // 올림 나눗셈
		m:    m,
		k:    k,
	}
}

// NewWithEstimates는 예상 원소 수 n과 목표 false positive 확률 p로부터
// m과 k를 자동 계산해 생성한다.
func NewWithEstimates(n uint64, p float64) *BloomFilter {
	m := OptimalM(n, p)
	return New(m, OptimalK(m, n))
}

// Cap은 전체 비트 수 m을 반환한다.
func (f *BloomFilter) Cap() uint64 { return f.m }

// K는 해시 함수 개수를 반환한다.
func (f *BloomFilter) K() uint64 { return f.k }

// Count는 지금까지 추가된 원소 수를 반환한다.
func (f *BloomFilter) Count() uint64 { return f.n }

// OptimalM은 원소 수 n과 목표 false positive 확률 p에 대해 필요한 비트 수를 계산한다.
// p가 (0, 1) 범위를 벗어나면 유효하지 않은 값으로 간주해 기본값 0.01로 대체한다.
//
//	m = -n * ln(p) / (ln2)^2
func OptimalM(n uint64, p float64) uint64 {
	if n == 0 {
		return 1
	}
	if p <= 0 || p >= 1 {
		p = 0.01
	}
	m := -float64(n) * math.Log(p) / (math.Ln2 * math.Ln2)
	return uint64(math.Ceil(m))
}

// OptimalK는 비트 수 m과 원소 수 n에 대해 최적의 해시 함수 개수를 계산한다.
//
//	k = (m/n) * ln2
func OptimalK(m, n uint64) uint64 {
	if n == 0 {
		return 1
	}
	k := (float64(m) / float64(n)) * math.Ln2
	if k < 1 {
		return 1
	}
	return uint64(math.Round(k))
}

// hashes는 xxhash 64비트 결과 하나를 상·하위 32비트로 쪼개
// 이중 해싱(Kirsch-Mitzenmacher)에 쓸 두 값을 만든다.
// 해시 함수를 k개 따로 두지 않아도 통계적으로 동등한 분포를 얻는다.
//
// 주의: h2가 0이거나 짝수이면 인덱스가 일부 위치에 뭉치므로 홀수로 보정한다.
func (f *BloomFilter) hashes(data []byte) (uint64, uint64) {
	sum := xxhash.Sum64(data)
	h1 := sum & 0xffffffff
	h2 := (sum >> 32) | 1
	return h1, h2
}

// Add는 원소를 필터에 추가한다.
func (f *BloomFilter) Add(data []byte) {
	h1, h2 := f.hashes(data)
	for i := uint64(0); i < f.k; i++ {
		pos := (h1 + i*h2) % f.m
		f.bits[pos/64] |= 1 << (pos % 64)
	}
	f.n++
}

// Contains는 원소가 있을 수 있는지 확인한다.
// false면 확실히 없고, true면 있을 수도 있다.
func (f *BloomFilter) Contains(data []byte) bool {
	h1, h2 := f.hashes(data)
	for i := uint64(0); i < f.k; i++ {
		pos := (h1 + i*h2) % f.m
		if f.bits[pos/64]&(1<<(pos%64)) == 0 {
			return false
		}
	}
	return true
}
