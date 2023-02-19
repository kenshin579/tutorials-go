package go_miniredis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	s := miniredis.RunT(t)

	// Optionally set some keys your code expects:
	s.Set("foo", "bar")
	s.HSet("some", "other", "key")

	// Run your code and see if it behaves.
	// An example using the redigo library from "github.com/gomodule/redigo/redis":
	c, err := redis.Dial("tcp", s.Addr())
	assert.NoError(t, err)
	_, err = c.Do("SET", "foo", "bar")
	assert.NoError(t, err)

	// Optionally check values in redis...
	get, err := s.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", get)

	// ... or use a helper for that:
	s.CheckGet(t, "foo", "bar")

	// TTL and expiration:
	err = s.Set("foo", "bar")
	assert.NoError(t, err)

	s.SetTTL("foo", 10*time.Second)
	s.FastForward(11 * time.Second)
	assert.False(t, s.Exists("foo"))
}
