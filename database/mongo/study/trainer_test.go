package study

import (
	"context"
	"testing"

	"github.com/chidiwilliams/flatbson"
	"github.com/kenshin579/tutorials-go/database/mongo/config"
	"github.com/kenshin579/tutorials-go/database/mongo/domain"
	"github.com/kenshin579/tutorials-go/database/mongo/test/mongodb"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type trainerTestSuite struct {
	suite.Suite
	ctx        context.Context
	db         *mongo.Database
	collection string
}

func TestMongoSuite(t *testing.T) {
	suite.Run(t, new(trainerTestSuite))
}

func (suite *trainerTestSuite) SetupSuite() {
	ctx := context.Background()
	cfg, err := config.New("../test/config.yaml")
	suite.NoError(err)

	db, err := mongodb.NewMongoDB(cfg)
	suite.NoError(err)
	suite.ctx = ctx
	suite.db = db
	suite.collection = string(domain.CollectionTrainer)
}

func (suite *trainerTestSuite) TestTrainer_Create() {
	suite.Run("create trainer", func() {
		// GIVEN
		sample := domain.CreateTrainerSample(1)

		// WHEN
		result, err := suite.createTrainer(suite.ctx, sample)
		defer suite.deleteTrainer(suite.ctx, sample.ID)

		// THEN
		suite.NoError(err)
		suite.Equal(sample.ID, result.ID)

	})
}

func (suite *trainerTestSuite) TestTrainer_Get() {
	suite.Run("create trainer", func() {
		// GIVEN
		sample := domain.CreateTrainerSample(1)
		_, err := suite.createTrainer(suite.ctx, sample)
		suite.NoError(err)
		defer suite.deleteTrainer(suite.ctx, sample.ID)

		// WHEN
		result, err := suite.getTrainer(suite.ctx, sample.ID)

		// THEN
		suite.NoError(err)
		suite.Equal(sample.ID, result.ID)
	})
}

func (suite *trainerTestSuite) TestTrainer_Update() {
	suite.Run("create trainer", func() {
		// GIVEN
		sample := domain.CreateTrainerSample(1)
		suite.createTrainer(suite.ctx, sample)
		defer suite.deleteTrainer(suite.ctx, sample.ID)
		updateRequest := domain.Trainer{
			ID:   sample.ID,
			Name: "Frank",
		}

		// WHEN
		result, err := suite.updateTrainer(suite.ctx, updateRequest)

		// THEN
		suite.NoError(err)
		suite.Equal(sample.ID, result.ID)
		suite.Equal(updateRequest.Name, result.Name)
	})
}

func (suite *trainerTestSuite) TestTrainer_Delete() {
	suite.Run("create trainer", func() {
		// GIVEN
		sample := domain.CreateTrainerSample(1)
		_, err := suite.createTrainer(suite.ctx, sample)
		suite.NoError(err)

		// WHEN
		err = suite.deleteTrainer(suite.ctx, sample.ID)

		// THEN
		suite.NoError(err)
		_, err = suite.getTrainer(suite.ctx, sample.ID)
		suite.Error(err, "mongo: no documents in result")
	})
}

func (suite *trainerTestSuite) createTrainer(ctx context.Context, trainer domain.Trainer) (domain.Trainer, error) {
	_, err := suite.db.Collection(suite.collection).InsertOne(ctx, trainer)
	if err != nil {
		return domain.Trainer{}, err
	}
	return trainer, nil
}

func (suite *trainerTestSuite) getTrainer(ctx context.Context, trainerID string) (domain.Trainer, error) {
	var trainer domain.Trainer
	if err := suite.db.Collection(suite.collection).FindOne(ctx, bson.M{"_id": trainerID}).Decode(&trainer); err != nil {
		return domain.Trainer{}, err
	}
	return trainer, nil
}

func (suite *trainerTestSuite) updateTrainer(ctx context.Context, trainer domain.Trainer) (domain.Trainer, error) {
	updateTrainer, err := flatbson.Flatten(trainer)
	if err != nil {
		return domain.Trainer{}, err
	}

	if err := suite.db.Collection(suite.collection).FindOneAndUpdate(ctx,
		bson.D{{Key: "_id", Value: trainer.ID}},
		bson.D{{Key: "$set", Value: updateTrainer}}).Err(); err != nil {
		return domain.Trainer{}, err
	}

	updatedTrainer, err := suite.getTrainer(ctx, trainer.ID)
	if err != nil {
		return domain.Trainer{}, err
	}
	return updatedTrainer, nil
}

func (suite *trainerTestSuite) deleteTrainer(ctx context.Context, trainerID string) error {
	_, err := suite.db.Collection(suite.collection).DeleteOne(ctx, bson.D{{
		Key:   "_id",
		Value: trainerID,
	}})
	if err != nil {
		return err
	}
	return nil
}
