package inmemory

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func NewRedisDB() (*miniredis.Miniredis, redis.UniversalClient) {
	miniRedis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: miniRedis.Addr()})

	return miniRedis, redisClient
}
