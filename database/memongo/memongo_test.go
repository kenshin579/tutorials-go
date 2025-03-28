package memongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/kenshin579/tutorials-go/database/memongo/db"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTestSuite struct {
	suite.Suite
	db     *mongo.Database
	ctx    context.Context
	rating *mongo.Collection
}

func TestMongoTestSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}

func (suite *MongoTestSuite) SetupSuite() {
	// create mongodb
	suite.db = db.NewInMemoryMongoDB()
	suite.rating = suite.db.Collection("coll_rating")
	suite.ctx = context.TODO()

	docs := []interface{}{
		bson.D{{"type", "Masala"}, {"rating", 10}, {"vendor", bson.A{"A", "C"}}},
		bson.D{{"type", "English Breakfast"}, {"rating", 6}},
		bson.D{{"type", "Oolong"}, {"rating", 7}, {"vendor", bson.A{"C"}}},
		bson.D{{"type", "Assam"}, {"rating", 5}},
		bson.D{{"type", "Earl Grey"}, {"rating", 8}, {"vendor", bson.A{"A", "B"}}},
	}
	_, err := suite.rating.InsertMany(suite.ctx, docs)
	require.NoError(suite.T(), err)
}

func (suite *MongoTestSuite) SetupTest() {
	fmt.Println("SetupTest")
}

func (suite *MongoTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite")
}

func (suite *MongoTestSuite) TestQuery() {
	ctx := suite.ctx

	filter := bson.D{{"type", "Oolong"}}
	cursor, err := suite.rating.Find(ctx, filter)

	var results []bson.D
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
