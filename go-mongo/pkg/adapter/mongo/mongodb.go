package mongo

import (
	"context"
	"log"

	"github.com/kenshin579/tutorials-go/go-mongo/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongodb struct {
	client *mongo.Client
}

func New(ctx context.Context, cfg *config.Config) (*Mongodb, error) {
	//client 옵션
	clientOptions := options.Client().ApplyURI(cfg.MongoConfig.Uri) /*.SetAuth(options.Credential{
		Username: cfg.MongoConfig.Username,
		Password: cfg.MongoConfig.Password,
	})*/

	//mongo 연결
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Mongodb{
		client: client,
	}, nil
}
