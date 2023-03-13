package go_testcontainers

import (
	"context"
	"testing"
	"time"

	redislib_v9 "github.com/go-redis/redis/v9"
	"github.com/kenshin579/tutorials-go/test/testcontainers"
	"github.com/stretchr/testify/suite"
)

type redisTestContainerTestSuite struct {
	suite.Suite
	redisV9Client *redislib_v9.Client
	ctx           context.Context
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(redisTestContainerTestSuite))
}

func (suite *redisTestContainerTestSuite) SetupSuite() {
	redisV9Client := testcontainers.NewRedisV9Client()
	suite.ctx = context.Background()
	suite.redisV9Client = redisV9Client

}

func (suite *redisTestContainerTestSuite) TearDownTest() {
	suite.NoError(suite.redisV9Client.FlushAll(suite.ctx).Err())
}

func (suite *redisTestContainerTestSuite) Test_RedisTestContainers() {
	suite.Run("Test Set and Get", func() {
		suite.redisV9Client.Set(suite.ctx, "hello", "world", time.Duration(0))

		value := suite.redisV9Client.Get(suite.ctx, "hello").Val()
		suite.Equal("world", value)
	})
}
