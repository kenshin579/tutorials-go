package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/kenshin579/tutorials-go/go-mongo/pkg/mongodb"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kenshin579/tutorials-go/go-mongo/pkg/config"
)

const (
	CollectionName = "trainers"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	cfg, err := config.New("go-mongo/config/config.yaml")

	if err != nil {
		log.Panic("config error: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := connect(ctx, cfg)

	ping(client, ctx)

	//define collection
	dbCollection := client.Database(cfg.MongoDBConfig.Database).Collection(CollectionName)

	insertOne(dbCollection, ctx)

	insertMany(dbCollection, ctx)

	update(dbCollection, ctx)

	findOne(dbCollection, ctx)

	findMultipleDocuments(dbCollection, ctx)

	deleteMany(client, dbCollection, ctx)
}

func deleteMany(client *mongo.Client, dbCollection *mongo.Collection, ctx context.Context) {
	//Delete all the documents in the dbCollection
	deleteResult, err := dbCollection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted %v documents in the trainers dbCollection\n", deleteResult.DeletedCount)

	// Close the connection once no longer needed
	err = client.Disconnect(ctx)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}
}

func findMultipleDocuments(dbCollection *mongo.Collection, ctx context.Context) {
	//findOptions을 주지 않고 검색하면 limit없이 검색이 된다
	findOptions := options.Find()
	findOptions.SetLimit(2) //최대 검색 객수 2개로 제한함

	var results []*Trainer

	// Finding multiple documents returns a cursor
	cur, err := dbCollection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the cursor
	for cur.Next(ctx) {
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(ctx)

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", len(results))

	for _, result := range results {
		fmt.Println(result)
	}
}

func findOne(dbCollection *mongo.Collection, ctx context.Context) {
	// Find a single document
	var result Trainer
	filter := bson.D{{"name", "Ash"}}

	err := dbCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)
}

func update(dbCollection *mongo.Collection, ctx context.Context) {
	// Update a document
	filter := bson.D{{"name", "Ash"}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	updateResult, err := dbCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func insertMany(c *mongo.Collection, ctx context.Context) {

	// Insert multiple documents
	data := getInitData()
	trainers := []interface{}{data[1], data[2]}

	insertManyResult, err := c.InsertMany(ctx, trainers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

func insertOne(c *mongo.Collection, ctx context.Context) {
	// Get a handle for your collection

	trainers := getInitData()

	// Insert a single document
	insertResult, err := c.InsertOne(ctx, trainers[0])

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func getInitData() []Trainer {
	var trainers []Trainer

	// Some dummy data to add to the Database
	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}
	trainers = append(trainers, ash)
	trainers = append(trainers, misty)
	trainers = append(trainers, brock)
	return trainers
}

func ping(c *mongo.Client, ctx context.Context) {
	// Check the connection
	err := c.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

func connect(ctx context.Context, cfg *config.Config) *mongo.Client {
	m, _ := mongodb.New(ctx, cfg)
	return m.Client
}
