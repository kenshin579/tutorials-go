package domain

import "github.com/go-redis/redis"

type RedisClient interface {
	redis.UniversalClient
}
