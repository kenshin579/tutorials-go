package store

import (
	"context"
	"fmt"
	"log"

	"github.com/kenshin579/tutorials-go/database/mongo/adapter/mongodb"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/kenshin579/tutorials-go/database/mongo/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStore struct {
	db *mongo.Database
}

func NewMongoStore(db *mongodb.Mongodb) *mongoStore {
	return &mongoStore{
		db: db.DB,
	}
}

func (m *mongoStore) Insert(ctx context.Context, trainer domain.Trainer) error {
	dbCollection := m.db.Collection(domain.CollectionName)

	result, err := dbCollection.InsertOne(ctx, trainer)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Inserted a single document: ", result.InsertedID)
	return nil
}

func (m *mongoStore) FindOne(ctx context.Context, filter interface{}) (domain.Trainer, error) {
	dbCollection := m.db.Collection(domain.CollectionName)

	// Find a single document
	var result domain.Trainer

	err := dbCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
		return result, err
	}

	fmt.Printf("Found a single document: %+v\n", result)
	return result, nil
}

func (m *mongoStore) InsertMany(ctx context.Context, trainerList []domain.Trainer) error {
	dbCollection := m.db.Collection(domain.CollectionName)
	var trainers []interface{}

	for _, t := range trainerList {
		trainers = append(trainers, t)
	}

	insertManyResult, err := dbCollection.InsertMany(ctx, trainers)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	return nil
}

func (m *mongoStore) FindAll(ctx context.Context, findOptions *options.FindOptions) ([]*domain.Trainer, error) {
	dbCollection := m.db.Collection(domain.CollectionName)
	// findOptions을 주지 않고 검색하면 limit없이 검색이 된다
	// findOptions := options.Find()
	// //findOptions.SetLimit(2) //최대 검색 객수 2개로 제한함

	var result []*domain.Trainer

	// Finding multiple documents returns a cursor
	cur, err := dbCollection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Iterate through the cursor
	for cur.Next(ctx) {
		var trainer domain.Trainer
		err := cur.Decode(&trainer)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, &trainer)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Close the cursor once finished
	cur.Close(ctx)

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", len(result))
	return result, nil
}

func (m *mongoStore) Update(ctx context.Context, filter interface{}, update interface{}) error {
	dbCollection := m.db.Collection(domain.CollectionName)

	updateResult, err := dbCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return nil
}

func (m *mongoStore) Delete(ctx context.Context, filter interface{}) error {
	dbCollection := m.db.Collection(domain.CollectionName)

	deleteResult, err := dbCollection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Deleted %v documents in the trainers dbCollection\n", deleteResult.DeletedCount)
	return nil
}
