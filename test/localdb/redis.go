package localdb

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewRedisClient() *redis.Client {
	endPoint, err := createMysqlTestContainer()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: endPoint,
	})

	return client
}

func createMysqlTestContainer() (string, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:6",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	endPoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		panic(err)
	}
	return endPoint, nil
}
