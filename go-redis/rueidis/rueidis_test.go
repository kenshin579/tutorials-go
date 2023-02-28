package rueidis

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/suite"
)

type rueidisTestSuite struct {
	suite.Suite
	ctx    context.Context
	client redis.UniversalClient
}

func TestRueidisSuite(t *testing.T) {
	suite.Run(t, new(rueidisTestSuite))
}

func (suite *rueidisTestSuite) SetupSuite() {
	fmt.Println("TestRueidisSuite started")

	suite.ctx = context.Background()
	suite.client = newRedisClient()

}

func (suite *rueidisTestSuite) TearDownTest() {
	// client.FlushAll(context.Background())
	_ = suite.client.Close()
}

func (suite *rueidisTestSuite) TestReJSON() {
	suite.Run("Ping", func() {
		pong, err := suite.client.Ping(context.Background()).Result()

		suite.NoError(err)
		suite.Equal("PONG", pong)
	})

	suite.Run("Set and Get Test", func() {
		// student := domain.Student{
		// 	Name: domain.Name{
		// 		First:  "Mark",
		// 		Middle: "S",
		// 		Last:   "Pronto",
		// 	},
		// 	Rank: 1,
		// }

	})

}

func newRedisClient() redis.UniversalClient {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return c
}
