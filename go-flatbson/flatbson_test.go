package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID      string  `bson:"_id,omitempty"`
	Name    string  `bson:"name,omitempty"`
	Address Address `bson:"address,omitempty"`
}

type Address struct {
	Street    string    `bson:"street,omitempty"`
	City      string    `bson:"city,omitempty"`
	State     string    `bson:"state,omitempty"`
	VisitedAt time.Time `bson:"visitedAt,omitempty"`
}

const TestCollection = ""

func Example() {
	//flatbson.Flatten(User{Address: {VisitedAt: time.Now().UTC()}}) //todo: 왜 오류가 발생하나?

	//flatbson.Flatten(User{Address: {Street: ""}})
	client := connect()
	client.Database(TestCollection)
}

func connect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/go-mongo-test").SetAuth(
		options.Credential{
			Username: "root",
			Password: "password",
		})

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
