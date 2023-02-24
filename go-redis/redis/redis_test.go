package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kenshin579/tutorials-go/go-redis/model"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var (
	client *redis.Client
)

func setup() {
	client = newRedisClient()
}

func teardown() {
	// client.FlushAll(context.Background())
	client.Close()
}

func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func Test_Ping(t *testing.T) {
	setup()
	defer teardown()

	pong, err := client.Ping(context.Background()).Result()

	assert.NoError(t, err)
	assert.Equal(t, "PONG", pong)
}

func Test_Set_Get_With_Primitive_Data_Type(t *testing.T) {
	setup()
	defer teardown()

	const TestValue = "Elliot"

	err := client.Set(context.Background(), "name", TestValue, 0).Err()

	assert.NoError(t, err)

	val, err := client.Get(context.Background(), "name").Result()
	assert.NoError(t, err)
	assert.Equal(t, TestValue, val)
}

func Test_Set_Get_With_Struct(t *testing.T) {
	setup()
	defer teardown()

	const TestKey = "id1234"

	authorJson, err := json.Marshal(model.Author{Name: "Elliot", Age: 25})
	assert.NoError(t, err)
	err = client.Set(context.Background(), TestKey, authorJson, 0).Err()
	val, err := client.Get(context.Background(), TestKey).Result()

	fmt.Printf("%v %T\n", val, val)
	var a model.Author

	err = json.Unmarshal([]byte(val), &a)
	assert.NoError(t, err)

	assert.Equal(t, "Elliot", a.Name)
}
