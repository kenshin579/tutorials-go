package counter

import (
	"sync"
	"sync/atomic"
)

// Counter Stores counts associated with a key.
type Counter struct {
	m sync.Map
}

// Get Retrieves the count without modifying it
func (c *Counter) Get(key string) (int64, bool) {
	count, ok := c.m.Load(key)
	if ok {
		return atomic.LoadInt64(count.(*int64)), true
	}
	return 0, false
}

// Add Adds value to the stored underlying value if it exists.
// If it does not exists, the value is assigned to the key.
func (c *Counter) Add(key string, value int64) int64 {
	count, loaded := c.m.LoadOrStore(key, &value)
	if loaded {
		return atomic.AddInt64(count.(*int64), value)
	}
	return *count.(*int64)
}

// DeleteAndGetLastValue Deletes the value associated with key and retrieves its.
func (c *Counter) DeleteAndGetLastValue(key string) (int64, bool) {
	lastValue, loaded := c.m.LoadAndDelete(key)
	if loaded {
		return *lastValue.(*int64), loaded
	}

	return 0, false
}
