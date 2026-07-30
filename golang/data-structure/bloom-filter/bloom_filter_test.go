package bloomfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptimalM(t *testing.T) {
	// n=100만, p=0.01 -> m = -n*ln(p) / (ln2)^2 = 9,585,059 비트 (약 1.14MiB)
	assert.Equal(t, uint64(9585059), OptimalM(1_000_000, 0.01))

	// p가 작아질수록 더 많은 비트가 필요하다
	assert.Greater(t, OptimalM(1_000_000, 0.001), OptimalM(1_000_000, 0.01))
}

func TestOptimalK(t *testing.T) {
	// k = (m/n)*ln2 = 9.585059 * 0.693147 = 6.64 -> 반올림 7
	assert.Equal(t, uint64(7), OptimalK(9585059, 1_000_000))

	// 최소 1개는 보장한다
	assert.Equal(t, uint64(1), OptimalK(10, 1_000_000))
}
