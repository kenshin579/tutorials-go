package clock

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

// 패턴 ③(가짜 시계 주입)과 달리, synctest는 코드 수정 없이
// 실제 시계(RealClock)를 쓰는 캐시를 가상 시간 버블 안에서 테스트한다.
// 버블 안에서는 time.Now()/time.Sleep()이 가상 시간으로 동작한다 (Go 1.25+).
func Test_TTLCache_Synctest_실제_시계로_만료_검증(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cache := NewTTLCache(clockwork.NewRealClock(), 10*time.Minute)

		cache.Set("session", "user-42")

		// 버블 안의 time.Sleep은 실제로 기다리지 않고 가상 시간을 진행시킨다
		time.Sleep(10*time.Minute - time.Second)
		v, ok := cache.Get("session")
		assert.True(t, ok)
		assert.Equal(t, "user-42", v)

		time.Sleep(time.Second)
		_, ok = cache.Get("session")
		assert.False(t, ok)
	})
}
