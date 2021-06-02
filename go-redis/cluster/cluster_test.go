package cluster

import (
	"context"
	"testing"

	"github.com/kenshin579/tutorials-go/go-redis/cluster/config"

	"github.com/go-redis/redis/v8"
)

var (
	cfg *config.Config
)

func setup() {
	cfg, _ = config.New("config/config.yaml")
}

func Test_Ping2(t *testing.T) {
	setup()

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: cfg.RedisConfig.ServerList,
	})
	err := rdb.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})

	if err != nil {
		panic(err)
	}
}
