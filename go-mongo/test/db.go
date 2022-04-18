package test

import (
	"context"

	"github.com/benweissmann/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connect struct{}

func (Connect) NewDB(ctx context.Context) (*mongo.Database, error) {
	opts := &memongo.Options{
		MongoVersion:     "4.2.1",
		ShouldUseReplica: false,
	}

	mongoServer, err := memongo.StartWithOptions(opts)
	if err != nil {
		return nil, err
	}

	clientOpts := options.Client().ApplyURI(mongoServer.URIWithRandomDB())
	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	return mongoClient.Database(memongo.RandomDatabase()), nil
}
