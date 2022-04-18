package main

import (
	"context"

	"github.com/kenshin579/tutorials-go/go-mongo/test"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTestSuite struct {
	suite.Suite
	test.Connect
	db *mongo.Database
}

func (suite *mongoTestSuite) SetupSuite() {
	db, err := suite.NewDB(context.TODO())
	suite.NoError(err)

	suite.db = db
}

func (suite *mongoTestSuite) Test() {
	suite.Run("test", func() {

	})
}
