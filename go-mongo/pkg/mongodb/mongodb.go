package mongodb

import (
	"context"
	"log"

	"github.com/kenshin579/tutorials-go/go-mongo/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodb struct {
	Client *mongo.Client
}

func New(ctx context.Context, cfg *config.Config) (*mongodb, error) {
	//Client 옵션
	clientOptions := options.Client().ApplyURI(cfg.MongoDBConfig.Uri) /*.SetAuth(options.Credential{
		Username: cfg.MongoDBConfig.Username,
		Password: cfg.MongoDBConfig.Password,
	})*/

	//mongo 연결
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &mongodb{
		Client: client,
	}, nil
}
