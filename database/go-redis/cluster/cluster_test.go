package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kenshin579/tutorials-go/go-redis/model"

	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/tutorials-go/go-redis/cluster/config"

	"github.com/go-redis/redis/v8"
)

var (
	cfg           *config.Config
	clusterClient *redis.ClusterClient
)

func setup() {
	cfg, _ = config.New("config/config.yaml")
	clusterClient = newClusterClient()
}

func teardown() {
	clusterClient.Close()
}
func newClusterClient() *redis.ClusterClient {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cfg.RedisConfig.ClusterConfig.ServerList,
		Password: cfg.RedisConfig.ClusterConfig.Password,
	})
	return rdb
}

func Test_Ping(t *testing.T) {
	setup()
	defer teardown()

	clusterClient := newClusterClient()
	err := clusterClient.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})

	assert.NoError(t, err)

}

func Test_Set_Get_With_Primitive_Data_Type(t *testing.T) {
	setup()
	defer teardown()
	ctx := context.Background()

	const TestValue = "Elliot"

	err := clusterClient.Set(ctx, "name", TestValue, 0).Err()
	assert.NoError(t, err)

	val, err := clusterClient.Get(ctx, "name").Result()
	assert.NoError(t, err)
	assert.Equal(t, TestValue, val)

}

func Test_Set_Get_With_Struct(t *testing.T) {
	setup()
	defer teardown()
	ctx := context.Background()

	const TestKey = "id1234"

	authorJson, err := json.Marshal(model.Author{Name: "Elliot", Age: 25})
	assert.NoError(t, err)
	err = clusterClient.Set(ctx, TestKey, authorJson, 0).Err()
	val, err := clusterClient.Get(ctx, TestKey).Result()

	fmt.Printf("%v %T\n", val, val)
	var a model.Author

	err = json.Unmarshal([]byte(val), &a)
	assert.NoError(t, err)

	assert.Equal(t, "Elliot", a.Name)
}
