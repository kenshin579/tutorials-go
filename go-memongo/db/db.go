package db

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewInMemoryMongoDB() *mongo.Database {
	opts := &memongo.Options{
		MongoVersion:     "4.2.1",
		ShouldUseReplica: true,
	}

	mongoServer, err := memongo.StartWithOptions(opts)
	if err != nil {
		log.Fatalf(err.Error())
	}

	clientOpts := options.Client().ApplyURI(mongoServer.URI())
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return client.Database(memongo.RandomDatabase())
}
