package func_opts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	maxAge := 1 * time.Hour
	maxSize := 10000
	cache, _ := NewCache(CacheMaxAge(maxAge), CacheMaxEntries(maxSize))
	assert.Equal(t, maxSize, cache.MaxSize)
	assert.Equal(t, maxAge, cache.MaxAge)
}
