package mongodb

import (
	"context"
	"log"

	"github.com/kenshin579/tutorials-go/go-mongo/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(cfg *config.Config) (*mongo.Database, error) {
	ctx := context.Background()

	//client 옵션
	clientOptions := options.Client().ApplyURI(cfg.MongoConfig.Uri)

	//mongo 연결
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client.Database(cfg.MongoConfig.Database), nil
}
