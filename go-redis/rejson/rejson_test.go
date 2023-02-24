package rejson

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/kenshin579/tutorials-go/go-redis/domain"
	"github.com/nitishm/go-rejson/v4"
	"github.com/stretchr/testify/suite"
)

type rejsonTestSuite struct {
	suite.Suite
	ctx    context.Context
	client redis.UniversalClient
	rh     *rejson.Handler
}

func TestReJSONSuite(t *testing.T) {
	suite.Run(t, new(rejsonTestSuite))
}

func (suite *rejsonTestSuite) SetupSuite() {
	fmt.Println("TestReJSONSuite started")

	suite.ctx = context.Background()
	suite.client = newRedisClient()

	suite.rh = rejson.NewReJSONHandler()
	suite.rh.SetGoRedisClientWithContext(suite.ctx, suite.client)
}

func (suite *rejsonTestSuite) TearDownTest() {
	// client.FlushAll(context.Background())
	_ = suite.client.Close()
}

func (suite *rejsonTestSuite) TestReJSON() {
	suite.Run("Ping", func() {
		pong, err := suite.client.Ping(context.Background()).Result()

		suite.NoError(err)
		suite.Equal("PONG", pong)
	})

	suite.Run("Set and Get Test", func() {
		student := domain.Student{
			Name: domain.Name{
				First:  "Mark",
				Middle: "S",
				Last:   "Pronto",
			},
			Rank: 1,
		}
		result, err := suite.rh.JSONSet("student", ".", student)
		suite.NoError(err)
		suite.Equal(result.(string), "OK")

		result, err = redigo.Bytes(suite.rh.JSONGet("student", "."))
		suite.NoError(err)

		res, err := suite.rh.JSONGet("student", ".")
		suite.NoError(err)

		getStudent := domain.Student{}
		err = json.Unmarshal(res.([]byte), &getStudent)
		suite.NoError(err)
		suite.Equal(student, getStudent)

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
