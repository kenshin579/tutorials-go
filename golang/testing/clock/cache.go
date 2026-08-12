package clock

import (
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
)

type cacheItem struct {
	value     string
	expiresAt time.Time
}

// TTLCache는 TTL이 지나면 항목이 만료되는 인메모리 캐시다.
// clockwork.Clock을 주입받아 테스트에서 가짜 시계로 시간을 "진행"시킬 수 있다.
// nowFunc 주입(패턴 ②)과 달리 시간이 흘러야 검증되는 로직에 적합하다.
type TTLCache struct {
	clock clockwork.Clock
	ttl   time.Duration

	mu    sync.Mutex
	items map[string]cacheItem
}

// NewTTLCache는 주어진 시계와 TTL로 캐시를 생성한다.
// 프로덕션에서는 clockwork.NewRealClock()을 전달한다.
func NewTTLCache(clock clockwork.Clock, ttl time.Duration) *TTLCache {
	return &TTLCache{
		clock: clock,
		ttl:   ttl,
		items: make(map[string]cacheItem),
	}
}

// Set은 key에 value를 저장하고 TTL 이후 만료되도록 기록한다.
func (c *TTLCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		value:     value,
		expiresAt: c.clock.Now().Add(c.ttl),
	}
}

// Get은 key의 값을 반환한다. 항목이 없거나 만료됐으면 false를 반환한다.
func (c *TTLCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	it, ok := c.items[key]
	if !ok {
		return "", false
	}
	// now >= expiresAt이면 만료
	if !c.clock.Now().Before(it.expiresAt) {
		delete(c.items, key)
		return "", false
	}
	return it.value, true
}
