package test

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() *mongo.Client {
	endPoint, _ := startMongoContainer()

	uri := fmt.Sprintf("mongodb://%s", endPoint)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client

}

func startMongoContainer() (string, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:4.4.4-bionic",
		ExposedPorts: []string{"27017/tcp", "27018/tcp"},
		// Env: map[string]string{
		//	"MONGO_INITDB_ROOT_USERNAME": "root",
		//	"MONGO_INITDB_ROOT_PASSWORD": "example",
		// },
	}

	mongodbC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	endPoint, err := mongodbC.Endpoint(ctx, "")
	if err != nil {
		panic(err)
	}

	return endPoint, nil
}
