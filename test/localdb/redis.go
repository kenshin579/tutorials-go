package localdb

import (
	"context"

	redislib_v8 "github.com/go-redis/redis/v8"
	redislib_v9 "github.com/go-redis/redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewRedisV8Client() *redislib_v8.Client {
	endPoint, err := createMysqlTestContainer()
	if err != nil {
		panic(err)
	}

	client := redislib_v8.NewClient(&redislib_v8.Options{
		Addr: endPoint,
	})
	return client
}

func NewRedisV9Client() *redislib_v9.Client {
	endPoint, err := createMysqlTestContainer()
	if err != nil {
		panic(err)
	}

	client := redislib_v9.NewClient(&redislib_v9.Options{
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
