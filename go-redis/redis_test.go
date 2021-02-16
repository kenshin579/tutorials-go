package go_redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
)

func Test_Ping(t *testing.T) {
	fmt.Println("Go Redis Tutorial")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

}
