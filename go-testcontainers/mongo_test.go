package go_testcontainers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTestContainerTestSuite struct {
	suite.Suite
	ctx context.Context

	mongoClient    *mongo.Client
	lockCollection *mongo.Collection
}

func TestMongoTestSuite(t *testing.T) {
	suite.Run(t, new(mongoTestContainerTestSuite))
}

func (suite *mongoTestContainerTestSuite) SetupSuite() {
	suite.ctx = context.Background()

}

func (suite *mongoTestContainerTestSuite) TearDownTest() {
	// suite.NoError(suite.redisV9Client.FlushAll(suite.ctx).Err())
}

func (suite *mongoTestContainerTestSuite) Test_MongoTestContainers() {
	suite.Run("Test Set and Get", func() {

	})
}
