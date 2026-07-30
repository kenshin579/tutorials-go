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

	// n이 0이면 비트가 필요 없지만 최소 1비트를 보장한다
	assert.Equal(t, uint64(1), OptimalM(0, 0.01))

	// p가 0 이하이면 유효하지 않으므로 기본값 0.01로 대체된다
	assert.Equal(t, OptimalM(1000, 0.01), OptimalM(1000, 0))

	// p가 1 이상이면 유효하지 않으므로 기본값 0.01로 대체된다
	assert.Equal(t, OptimalM(1000, 0.01), OptimalM(1000, 1.5))
}

func TestOptimalK(t *testing.T) {
	// k = (m/n)*ln2 = 9.585059 * 0.693147 = 6.64 -> 반올림 7
	assert.Equal(t, uint64(7), OptimalK(9585059, 1_000_000))

	// 최소 1개는 보장한다
	assert.Equal(t, uint64(1), OptimalK(10, 1_000_000))

	// n이 0이면 나눗셈이 불가능하므로 최소 1개를 보장한다
	assert.Equal(t, uint64(1), OptimalK(100, 0))
}

func TestNew(t *testing.T) {
	f := New(128, 3)

	assert.Equal(t, uint64(128), f.Cap())
	assert.Equal(t, uint64(3), f.K())
	assert.Equal(t, uint64(0), f.Count())
	// 128비트는 uint64 워드 2개
	assert.Len(t, f.bits, 2)
}

func TestNew_비트수가_64의_배수가_아닌_경우(t *testing.T) {
	f := New(100, 3)

	// 100비트를 담으려면 워드 2개가 필요하다
	assert.Len(t, f.bits, 2)
}

func TestNew_m이_0인_경우(t *testing.T) {
	// m이 0이면 방어 분기가 작동해 최소 1비트를 보장한다
	f := New(0, 3)

	assert.Equal(t, uint64(1), f.Cap())
}

func TestNew_k가_0인_경우(t *testing.T) {
	// k가 0이면 방어 분기가 작동해 최소 1개의 해시 함수를 보장한다
	f := New(128, 0)

	assert.Equal(t, uint64(1), f.K())
}

func TestNewWithEstimates(t *testing.T) {
	f := NewWithEstimates(1_000_000, 0.01)

	assert.Equal(t, uint64(9585059), f.Cap())
	assert.Equal(t, uint64(7), f.K())
}
