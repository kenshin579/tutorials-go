package bloomfilter

import (
	"fmt"
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/stretchr/testify/assert"
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

	// 직접 구현과 같은 m, k가 나오는지 확인한다
	assert.Equal(t, uint(9585059), filter.Cap())
	assert.Equal(t, uint(7), filter.K())
}

// TestAndAdd는 조회와 추가를 한 번에 처리한다. 중복 제거에 유용하다.
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

// 직접 구현과 라이브러리가 같은 원소에 대해 같은 판정을 내리는지 확인한다.
func TestLibrary_직접_구현과_동일하게_동작한다(t *testing.T) {
	mine := NewWithEstimates(10_000, 0.01)
	lib := bloom.NewWithEstimates(10_000, 0.01)

	for i := 0; i < 10_000; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		mine.Add(key)
		lib.Add(key)
	}

	// 넣은 원소는 양쪽 모두 true여야 한다 (false negative 없음)
	for i := 0; i < 10_000; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		assert.True(t, mine.Contains(key))
		assert.True(t, lib.Test(key))
	}
}
