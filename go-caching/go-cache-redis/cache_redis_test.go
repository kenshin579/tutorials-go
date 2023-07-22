package go_cache_redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/kenshin579/tutorials-go/test/inmemory"
	"github.com/stretchr/testify/suite"

	"github.com/go-redis/cache/v8"
)

type Object struct {
	Str string
	Num int
}

type cacheRedisTestSuite struct {
	suite.Suite
	ctx       context.Context
	miniredis *miniredis.Miniredis
}

func TestRedisStoreSuite(t *testing.T) {
	suite.Run(t, new(cacheRedisTestSuite))
}
func (suite *cacheRedisTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	db, _ := inmemory.NewRedisDB()
	suite.miniredis = db
}

func (suite *cacheRedisTestSuite) TearDownSuite() {
	suite.miniredis.Close()
}

func (suite *cacheRedisTestSuite) Test() {
	suite.Run("Ping", func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     suite.miniredis.Addr(),
			Password: "",
			DB:       0,
		})

		pong, err := rdb.Ping(suite.ctx).Result()
		suite.NoError(err)
		suite.Equal("PONG", pong)
	})

	suite.Run("Basic Usage", func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     suite.miniredis.Addr(),
			Password: "",
			DB:       0,
		})

		mycache := cache.New(&cache.Options{
			Redis:      rdb,
			LocalCache: cache.NewTinyLFU(10, time.Minute),
		})

		ctx := context.TODO()
		key := "mykey"
		obj := &Object{
			Str: "mystring",
			Num: 42,
		}

		if err := mycache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: obj,
			TTL:   time.Hour,
		}); err != nil {
			panic(err)
		}

		var wanted Object
		err := mycache.Get(ctx, key, &wanted)
		suite.NoError(err)
		suite.Equal(Object{Str: "mystring", Num: 42}, wanted)

	})

	suite.Run("Advanced Usage", func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     suite.miniredis.Addr(),
			Password: "",
			DB:       0,
		})

		mycache := cache.New(&cache.Options{
			Redis:      rdb,
			LocalCache: cache.NewTinyLFU(10, time.Minute),
		})

		obj := new(Object)
		err := mycache.Once(&cache.Item{
			Key:   "mykey",
			Value: obj, // destination
			Do: func(*cache.Item) (interface{}, error) {
				return &Object{
					Str: "mystring",
					Num: 42,
				}, nil
			},
		})
		suite.NoError(err)

		suite.Equal(&Object{Str: "mystring", Num: 42}, obj)
	})
}
