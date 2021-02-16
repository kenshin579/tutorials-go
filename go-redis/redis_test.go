package go_redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
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

func Test_Set_Value(t *testing.T) {
	const TestValue = "Elliot"

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := client.Set("name", TestValue, 0).Err()
	// if there has been an error setting the TestValue
	// handle the error
	if err != nil {
		fmt.Println(err)
	}

	val, err := client.Get("name").Result()
	if err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, TestValue, val)
}

type Author struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Test_json_struct(t *testing.T) {
	const TestKey = "id1234"

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	authorJson, err := json.Marshal(Author{Name: "Elliot", Age: 25})
	if err != nil {
		fmt.Println(err)
	}

	err = client.Set(TestKey, authorJson, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	val, err := client.Get(TestKey).Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v %T", val, val)
	var a Author

	json.Unmarshal([]byte(val), &a)

	assert.Equal(t, "Elliot", a.Name)
}
