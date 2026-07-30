// Package bloomfilter는 Bloom Filter를 비트 배열과 이중 해싱으로 직접 구현한 예제이다.
package bloomfilter

import "math"

// OptimalM은 원소 수 n과 목표 false positive 확률 p에 대해 필요한 비트 수를 계산한다.
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
