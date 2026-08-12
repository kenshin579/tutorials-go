package clock

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

func Test_TTLCache_가짜_시계로_만료_검증(t *testing.T) {
	fakeClock := clockwork.NewFakeClock()
	cache := NewTTLCache(fakeClock, 10*time.Minute)

	cache.Set("session", "user-42")

	// TTL 직전: 아직 살아있다
	fakeClock.Advance(10*time.Minute - time.Second)
	v, ok := cache.Get("session")
	assert.True(t, ok)
	assert.Equal(t, "user-42", v)

	// TTL 경과: 만료된다 (실제로 10분을 기다리지 않는다)
	fakeClock.Advance(time.Second)
	_, ok = cache.Get("session")
	assert.False(t, ok)
}

func Test_TTLCache_없는_키(t *testing.T) {
	cache := NewTTLCache(clockwork.NewFakeClock(), time.Minute)

	_, ok := cache.Get("missing")

	assert.False(t, ok)
}
