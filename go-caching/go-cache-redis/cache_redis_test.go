package go_cache_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/go-redis/cache/v8"
)

type Object struct {
	Str string
	Num int
}

func Example_basicUsage() {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"localhost": ":6379",
		},
	})

	mycache := cache.New(&cache.Options{
		Redis:      ring,
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
	if err := mycache.Get(ctx, key, &wanted); err == nil {
		fmt.Println(wanted)
	}

	// Output: {mystring 42}
}

func Example_advancedUsage() {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"localhost": ":6379",
		},
	})

	mycache := cache.New(&cache.Options{
		Redis:      ring,
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
	if err != nil {
		panic(err)
	}
	fmt.Println(obj)

	// Output: &{mystring 42}
}
