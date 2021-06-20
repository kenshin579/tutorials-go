package store

import (
	"fmt"
	"log"

	"github.com/kenshin579/tutorials-go/go-mongo/pkg/adapter/mongo"

	"github.com/kenshin579/tutorials-go/go-mongo/domain"
)

type mongoStore struct {
	Mongodb *mongo.Mongodb
}

func NewMongoStore(mongodb mongo.Mongodb) *mongoStore {
	return &mongoStore{mongoClient: client}
}

func (m *mongoStore) Insert(trainer domain.Trainer) {
	// Insert a single document
	insertResult, err := m.mongoClient.
		mongoClientInsertOne(ctx, trainers[0])

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func createSampleTrainer() []domain.Trainer {
	var trainers []domain.Trainer

	// Some dummy data to add to the Database
	ash := domain.Trainer{Name: "Ash", Age: 10, City: "Pallet Town"}
	misty := domain.Trainer{Name: "Misty", Age: 10, City: "Cerulean City"}
	brock := domain.Trainer{Name: "Brock", Age: 15, City: "Pewter City"}
	trainers = append(trainers, ash)
	trainers = append(trainers, misty)
	trainers = append(trainers, brock)
	return trainers
}
