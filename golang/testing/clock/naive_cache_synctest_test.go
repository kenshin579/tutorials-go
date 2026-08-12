package clock

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

// NaiveCache는 time.Now()를 직접 호출하고 시계 주입 설계가 없다.
// 그런데도 synctest 버블 안에서는 time.Now()/time.Sleep()이
// 가상 시간으로 동작하므로 코드 수정 없이 테스트할 수 있다 (Go 1.25+).
func Test_NaiveCache_Synctest_만료_검증(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cache := NewNaiveCache(10 * time.Minute)

		cache.Set("session", "user-42")

		// 버블 안의 time.Sleep은 실제로 기다리지 않고 가상 시간을 진행시킨다
		time.Sleep(10*time.Minute - time.Second)
		v, ok := cache.Get("session")
		assert.True(t, ok)
		assert.Equal(t, "user-42", v)

		// TTL 경과: 만료된다
		time.Sleep(time.Second)
		_, ok = cache.Get("session")
		assert.False(t, ok)
	})
}
