package go_testcontainers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTestContainerTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *mongo.Client
	collection *mongo.Collection
}

func TestMongoTestSuite(t *testing.T) {
	suite.Run(t, new(mongoTestContainerTestSuite))
}

func (s *mongoTestContainerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.client = NewMongoClient()
	s.collection = s.client.Database("testdb").Collection("items")
}

func (s *mongoTestContainerTestSuite) TearDownTest() {
	s.NoError(s.collection.Drop(s.ctx))
}

func (s *mongoTestContainerTestSuite) Test_InsertAndFind() {
	doc := bson.M{"name": "testcontainers", "language": "go"}
	_, err := s.collection.InsertOne(s.ctx, doc)
	s.NoError(err)

	var result bson.M
	err = s.collection.FindOne(s.ctx, bson.M{"name": "testcontainers"}).Decode(&result)
	s.NoError(err)
	s.Equal("testcontainers", result["name"])
	s.Equal("go", result["language"])
}

func (s *mongoTestContainerTestSuite) Test_InsertManyAndCount() {
	docs := []interface{}{
		bson.M{"name": "redis", "type": "cache"},
		bson.M{"name": "mongo", "type": "document"},
		bson.M{"name": "mysql", "type": "relational"},
	}
	_, err := s.collection.InsertMany(s.ctx, docs)
	s.NoError(err)

	count, err := s.collection.CountDocuments(s.ctx, bson.M{})
	s.NoError(err)
	s.Equal(int64(3), count)
}

func (s *mongoTestContainerTestSuite) Test_UpdateOne() {
	doc := bson.M{"name": "old", "version": 1}
	_, err := s.collection.InsertOne(s.ctx, doc)
	s.NoError(err)

	_, err = s.collection.UpdateOne(s.ctx,
		bson.M{"name": "old"},
		bson.M{"$set": bson.M{"name": "new", "version": 2}},
	)
	s.NoError(err)

	var result bson.M
	s.NoError(s.collection.FindOne(s.ctx, bson.M{"name": "new"}).Decode(&result))
	s.Equal(int32(2), result["version"])
}

func (s *mongoTestContainerTestSuite) Test_DeleteOne() {
	doc := bson.M{"name": "to-delete"}
	_, err := s.collection.InsertOne(s.ctx, doc)
	s.NoError(err)

	res, err := s.collection.DeleteOne(s.ctx, bson.M{"name": "to-delete"})
	s.NoError(err)
	s.Equal(int64(1), res.DeletedCount)

	count, err := s.collection.CountDocuments(s.ctx, bson.M{"name": "to-delete"})
	s.NoError(err)
	s.Equal(int64(0), count)
}
