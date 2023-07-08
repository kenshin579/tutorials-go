package go_miniredis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/kenshin579/tutorials-go/test"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type miniRedisTestSuite struct {
	suite.Suite
	ctx         context.Context
	rediServer  *miniredis.Miniredis
	redisClient redis.UniversalClient
}

func TestMiniRedisTestSuite(t *testing.T) {
	suite.Run(t, new(miniRedisTestSuite))
}

func (suite *miniRedisTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	server, client := test.NewRedisDB()
	suite.rediServer = server
	suite.redisClient = client

}

func (suite *miniRedisTestSuite) TearDownTest() {
	suite.rediServer.Close()
}

func (suite *miniRedisTestSuite) Test_MiniRedis() {
	suite.Run("Test Set and Get", func() {
		suite.redisClient.Set(suite.ctx, "hello", "world", time.Duration(0))
		suite.Equal("world", suite.redisClient.Get(suite.ctx, "hello").Val())
	})
}

func Test_Simple(t *testing.T) {
	s := miniredis.RunT(t)

	// Optionally set some keys your code expects:
	s.Set("foo", "bar")
	s.HSet("some", "other", "key")

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
