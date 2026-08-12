package clock

import (
	"sync"
	"time"
)

// NaiveCache는 time.Now()를 직접 호출하는 평범한 TTL 캐시다.
// 시계 주입 설계가 전혀 없어 패턴 1~3으로는 테스트할 수 없지만,
// testing/synctest 버블 안에서는 time.Now()가 가상 시간을 읽으므로
// 코드 수정 없이 그대로 테스트할 수 있다(패턴 4).
type NaiveCache struct {
	ttl time.Duration

	mu    sync.Mutex
	items map[string]cacheItem
}

// NewNaiveCache는 주어진 TTL로 캐시를 생성한다.
func NewNaiveCache(ttl time.Duration) *NaiveCache {
	return &NaiveCache{
		ttl:   ttl,
		items: make(map[string]cacheItem),
	}
}

// Set은 key에 value를 저장하고 TTL 이후 만료되도록 기록한다.
func (c *NaiveCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Get은 key의 값을 반환한다. 항목이 없거나 만료됐으면 false를 반환한다.
func (c *NaiveCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	it, ok := c.items[key]
	if !ok {
		return "", false
	}
	// now >= expiresAt이면 만료
	if !time.Now().Before(it.expiresAt) {
		delete(c.items, key)
		return "", false
	}
	return it.value, true
}
