package redis

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func Test_Ping(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()

	assert.NoError(t, err)
	assert.Equal(t, "PONG", pong)
}

func Test_Set_Get_With_Primitive_Data_Type(t *testing.T) {
	const TestValue = "Elliot"

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := client.Set("name", TestValue, 0).Err()
	assert.NoError(t, err)

	val, err := client.Get("name").Result()
	assert.NoError(t, err)
	assert.Equal(t, TestValue, val)
}

type Author struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Test_Set_Get_With_Struct(t *testing.T) {
	const TestKey = "id1234"

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	authorJson, err := json.Marshal(Author{Name: "Elliot", Age: 25})
	assert.NoError(t, err)
	err = client.Set(TestKey, authorJson, 0).Err()
	val, err := client.Get(TestKey).Result()

	fmt.Printf("%v %T\n", val, val)
	var a Author

	err = json.Unmarshal([]byte(val), &a)
	assert.NoError(t, err)

	assert.Equal(t, "Elliot", a.Name)
}
