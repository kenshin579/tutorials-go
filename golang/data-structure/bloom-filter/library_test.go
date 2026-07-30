package bloomfilter

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// bits-and-blooms/bloom 기본 사용법
func TestLibrary_기본_사용법(t *testing.T) {
	// 100만 개를 1% false positive 확률로 담는 필터
	filter := bloom.NewWithEstimates(1_000_000, 0.01)

	filter.Add([]byte("hello"))
	filter.AddString("world")

	assert.True(t, filter.Test([]byte("hello")))
	assert.True(t, filter.TestString("world"))
	assert.False(t, filter.TestString("없는-키"))

	// m은 양쪽이 항상 같다. k는 라이브러리가 올림(ceil), 직접 구현이 반올림(round)이라
	// p에 따라 1 차이가 날 수 있고, p=0.01에서는 우연히 일치한다.
	assert.Equal(t, OptimalM(1_000_000, 0.01), uint64(filter.Cap()))
	assert.Equal(t, uint(7), filter.K())
}

// TestAndAdd는 조회와 추가를 한 번에 처리한다. 중복 제거에 유용하다.
// 주의: false positive가 나면 처음 보는 항목을 이미 본 것으로 착각해 버린다. 유실이 허용되는 곳에만 쓴다.
func TestLibrary_TestAndAdd로_중복_제거(t *testing.T) {
	filter := bloom.NewWithEstimates(1000, 0.01)

	urls := []string{"a.com", "b.com", "a.com", "c.com", "b.com"}
	unique := 0
	for _, u := range urls {
		if !filter.TestAndAddString(u) {
			unique++
		}
	}

	assert.Equal(t, 3, unique)
}

// 두 구현은 실측 FPR이 비슷하지만, 해시가 다르므로
// 어떤 키에서 오답이 나는지는 서로 다르다.
func TestLibrary_실측FPR은_비슷하지만_오답_대상은_다르다(t *testing.T) {
	const (
		n      = 100_000
		trials = 100_000
	)

	mine := NewWithEstimates(n, 0.01)
	lib := bloom.NewWithEstimates(n, 0.01)
	for i := 0; i < n; i++ {
		key := []byte(fmt.Sprintf("member-%d", i))
		mine.Add(key)
		lib.Add(key)
	}

	var mineFP, libFP, bothFP int
	for i := 0; i < trials; i++ {
		key := []byte(fmt.Sprintf("stranger-%d", i))
		m, l := mine.Contains(key), lib.Test(key)
		if m {
			mineFP++
		}
		if l {
			libFP++
		}
		if m && l {
			bothFP++
		}
	}
	t.Logf("false positive: 직접 구현=%d 라이브러리=%d 둘 다=%d (조회 %d회)",
		mineFP, libFP, bothFP, trials)

	// 양쪽 모두 목표치 1% 근처다
	assert.InDelta(t, 0.01, float64(mineFP)/float64(trials), 0.003)
	assert.InDelta(t, 0.01, float64(libFP)/float64(trials), 0.003)

	// 해시가 독립적이므로 오답이 겹칠 기대값은 trials*0.01*0.01 = 10건 수준이다.
	// 양쪽이 같은 키에서 틀리는 일은 드물다.
	assert.Less(t, bothFP, mineFP/10)
}

// 직접 구현에는 없고 라이브러리에는 있는 기능 - 필터를 저장했다 복원할 수 있다.
func TestLibrary_직렬화_왕복(t *testing.T) {
	original := bloom.NewWithEstimates(1000, 0.01)
	original.AddString("hello")

	var buf bytes.Buffer
	_, err := original.WriteTo(&buf)
	require.NoError(t, err)

	restored := bloom.New(1, 1) // 크기는 ReadFrom이 덮어쓴다
	_, err = restored.ReadFrom(&buf)
	require.NoError(t, err)

	assert.True(t, restored.TestString("hello"))
	assert.Equal(t, original.Cap(), restored.Cap())
	assert.Equal(t, original.K(), restored.K())
}
