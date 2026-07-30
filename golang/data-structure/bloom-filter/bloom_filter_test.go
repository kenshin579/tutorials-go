package bloomfilter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		require.True(t, f.Contains([]byte(key)), "추가한 원소 %s 가 없다고 나왔다", key)
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

func TestBloomFilter_EstimatedFPR(t *testing.T) {
	f := NewWithEstimates(100_000, 0.01)

	// 비어 있으면 false positive가 날 수 없다
	assert.Equal(t, 0.0, f.EstimatedFPR())

	for i := 0; i < 50_000; i++ {
		f.Add([]byte(fmt.Sprintf("key-%d", i)))
	}
	half := f.EstimatedFPR()

	for i := 50_000; i < 100_000; i++ {
		f.Add([]byte(fmt.Sprintf("key-%d", i)))
	}
	full := f.EstimatedFPR()

	// 원소가 늘수록 false positive 확률은 커진다
	assert.Greater(t, full, half)
	// 설계 용량을 채웠을 때 목표치 1% 근처여야 한다.
	// 이 값은 난수가 개입하지 않는 순수 계산이므로 델타를 넉넉히 잡을 이유가 없다.
	assert.InDelta(t, 0.01, full, 0.001)
}

// 이론값과 실측값이 맞는지 확인한다. 본문 4.5절의 근거가 되는 테스트다.
func TestBloomFilter_실측_FalsePositiveRate(t *testing.T) {
	const (
		n      = 100_000   // 설계 용량
		trials = 1_000_000 // 넣지 않은 원소로 조회할 횟수
		target = 0.01      // 목표 false positive 확률
	)

	f := NewWithEstimates(n, target)
	for i := 0; i < n; i++ {
		f.Add([]byte(fmt.Sprintf("member-%d", i)))
	}

	falsePositives := 0
	for i := 0; i < trials; i++ {
		// "member-" 접두사와 겹치지 않는 키만 조회한다
		if f.Contains([]byte(fmt.Sprintf("stranger-%d", i))) {
			falsePositives++
		}
	}

	actual := float64(falsePositives) / float64(trials)
	t.Logf("m=%d k=%d n=%d", f.Cap(), f.K(), f.Count())
	t.Logf("false positive 건수=%d/%d", falsePositives, trials)
	t.Logf("이론 FPR=%.5f 실측 FPR=%.5f", f.EstimatedFPR(), actual)

	// 이론값과 실측값의 차이는 통계적 오차 범위 안에서만 나타나야 한다.
	// trials=100만, p=0.01이면 이항분포 표준편차가 약 99.5건(FPR로 1e-4)이므로
	// 델타 0.001은 약 10σ에 해당한다. 게다가 xxhash는 시드가 없고 입력이 고정되어
	// 매 실행 결과가 동일하므로 이 정도로 좁혀도 간헐적 실패가 나지 않는다.
	assert.InDelta(t, target, actual, 0.001)
}
