package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/kenshin579/tutorials-go/database/redis/model"

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
	client.Close()
}

// newRedisClient는 로컬 Redis에 접속하는 클라이언트를 생성한다.
// 비밀번호는 REDIS_PASSWORD 환경변수로 재정의할 수 있으며,
// 기본값은 cloud/docker/redis/Makefile의 requirepass 값과 동일하다.
func newRedisClient() *redis.Client {
	password := os.Getenv("REDIS_PASSWORD")
	if password == "" {
		password = "mypassword"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: password,
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
