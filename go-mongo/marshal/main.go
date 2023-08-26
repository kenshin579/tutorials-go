package main

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID       string     `bson:"_id"`
	Name     string     `bson:"name"`
	Birthday CustomTime `bson:"birthday"`
}

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	t := ct.Time
	return bson.MarshalValue(t)
}

func (ct *CustomTime) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var tTime time.Time
	err := bson.Unmarshal(data, &tTime)
	if err != nil {
		return err
	}
	ct.Time = tTime
	return nil
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx := context.Background()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("test2").Collection("people")

	p := Person{
		ID:       "p1",
		Name:     "John",
		Birthday: CustomTime{Time: time.Now()},
	}

	_, err = collection.InsertOne(ctx, p)
	if err != nil {
		log.Fatal(err)
	}

	// todo: error decoding key birthday: cannot decode invalid into a time.Time 문제가 있음
	var result Person
	err = collection.FindOne(ctx, bson.M{"_id": "p1"}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name: %s\n", result.Name)
	fmt.Printf("Birthday: %s\n", result.Birthday.Format("2006-01-02 15:04:05"))
}
