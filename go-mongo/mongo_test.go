package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/kenshin579/tutorials-go/go-mongo/config"
	"github.com/kenshin579/tutorials-go/go-mongo/domain"
	"github.com/kenshin579/tutorials-go/go-mongo/test/mongodb"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTestSuite struct {
	suite.Suite
	ctx context.Context
	db  *mongo.Database
}

func TestMongoSuite(t *testing.T) {
	suite.Run(t, new(mongoTestSuite))
}

//todo: 설정을 해야 함
func (suite *mongoTestSuite) SetupSuite() {
	ctx := context.Background()
	cfg, err := config.New("test/config.yaml")
	suite.NoError(err)

	db, err := mongodb.NewMongoDB(cfg)
	suite.NoError(err)
	suite.ctx = ctx
	suite.db = db
}

func (suite *mongoTestSuite) TestTrainer_Create() {
	suite.Run("create trainer", func() {
		trainer := domain.Trainer{
			Name: "Frank",
			Age:  20,
			City: "Columbus",
		}

		suite.createTrainer(suite.ctx, trainer)

	})
}

func (suite *mongoTestSuite) createTrainer(ctx context.Context, trainer domain.Trainer) {
	result, err := suite.db.Collection(string(domain.CollectionTrainer)).InsertOne(ctx, trainer)
	suite.NoError(err)
	fmt.Println(result)
}
