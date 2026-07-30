package bloomfilter

import (
	"fmt"
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

func TestBloomFilter_AddContains(t *testing.T) {
	f := NewWithEstimates(1000, 0.01)

	f.Add([]byte("hello"))
	f.Add([]byte("world"))

	assert.True(t, f.Contains([]byte("hello")))
	assert.True(t, f.Contains([]byte("world")))
	assert.Equal(t, uint64(2), f.Count())
}

func TestBloomFilter_빈_필터는_아무것도_포함하지_않는다(t *testing.T) {
	f := NewWithEstimates(1000, 0.01)

	assert.False(t, f.Contains([]byte("hello")))
}

// Bloom Filter의 핵심 보장: false negative는 절대 발생하지 않는다.
func TestBloomFilter_FalseNegative가_없다(t *testing.T) {
	const n = 100_000
	f := NewWithEstimates(n, 0.01)

	for i := 0; i < n; i++ {
		f.Add([]byte(fmt.Sprintf("member-%d", i)))
	}

	for i := 0; i < n; i++ {
		key := fmt.Sprintf("member-%d", i)
		assert.True(t, f.Contains([]byte(key)), "추가한 원소 %s 가 없다고 나왔다", key)
	}
}

// Add와 Contains는 힙 할당 없이 동작해야 한다.
// Task 5의 벤치마크에서 B/op를 라이브러리와 비교하려면 이 성질이 유지되어야 한다.
func TestBloomFilter_할당이_없다(t *testing.T) {
	f := NewWithEstimates(1000, 0.01)
	key := []byte("hello")

	addAllocs := testing.AllocsPerRun(100, func() {
		f.Add(key)
	})
	assert.Equal(t, 0.0, addAllocs)

	containsAllocs := testing.AllocsPerRun(100, func() {
		f.Contains(key)
	})
	assert.Equal(t, 0.0, containsAllocs)
}
