package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/kenshin579/tutorials-go/go-mongo/adapter/mongodb"

	"github.com/chidiwilliams/flatbson"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/kenshin579/tutorials-go/go-flatbson/model"

	"github.com/kenshin579/tutorials-go/go-mongo/config"

	"go.mongodb.org/mongo-driver/mongo"
)

const TestDatabaseName = "go-flatbson"

var TestCollectionName = getCurrentFilename()

func Test_Mongo_Update로_업데이트하는_경우_해당데이터_전체가_replace된다(t *testing.T) {
	cfg, err := config.New("../go-mongo/config/config.yaml")

	if err != nil {
		log.Panic("config error: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := connect(ctx, cfg)
	collection := client.Database(TestDatabaseName).Collection(TestCollectionName)
	setup(ctx, collection)
	defer teardown(ctx, collection)

	filter := bson.D{{"_id", "user-01"}}

	user := model.User{
		Address: model.Address{
			Street: "street2",
		},
	}
	update := bson.D{
		{"$set", user}}

	//Address 데이터 전테가 다 replace 되어버린다
	updateResult, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	assert.Equal(t, 1, int(updateResult.ModifiedCount))

	findUser := model.User{}
	collection.FindOne(ctx, bson.D{{"_id", "user-01"}}).Decode(&findUser)
	assert.Empty(t, findUser.Address.City)

}

func Test(t *testing.T) {
	user := model.User{
		ID:   "user-01",
		Name: "Frank",
		Address: model.Address{
			Street: "street1",
			City:   "seoul",
		},
	}

	flatten, err := flatbson.Flatten(user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v %T\n", flatten, flatten)
	assert.Equal(t, map[string]interface{}{
		"_id":            "user-01",
		"address.city":   "seoul",
		"address.street": "street1",
		"name":           "Frank",
	}, flatten)
}

func Test_Mongo_Update_With_FlatBson_Partially_업데이트가_된다(t *testing.T) {
	cfg, err := config.New("../go-mongo/config/config.yaml")

	if err != nil {
		log.Panic("config error: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := connect(ctx, cfg)
	collection := client.Database(TestDatabaseName).Collection(TestCollectionName)
	setup(ctx, collection)
	defer teardown(ctx, collection)

	filter := bson.D{{"_id", "user-01"}}

	user := model.User{
		Address: model.Address{
			Street: "street2",
		},
	}
	flatten, _ := flatbson.Flatten(user)
	update := bson.D{
		{"$set", flatten}}

	//Address 데이터 전테가 다 replace 되어버린다
	updateResult, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	assert.Equal(t, 1, int(updateResult.ModifiedCount))

	findUser := model.User{}
	collection.FindOne(ctx, bson.D{{"_id", "user-01"}}).Decode(&findUser)
	assert.Equal(t, "seoul", findUser.Address.City)
}

func setup(ctx context.Context, collection *mongo.Collection) {
	dateStr := "2020-07-03"
	parseDate, _ := time.Parse("2006-01-02", dateStr)

	user := model.User{
		ID:   "user-01",
		Name: "Frank",
		Address: model.Address{
			Street:    "street1",
			City:      "seoul",
			State:     "seoul",
			VisitedAt: parseDate,
		},
	}

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func teardown(ctx context.Context, collection *mongo.Collection) {
	collection.Drop(ctx)
}

func connect(ctx context.Context, cfg *config.Config) *mongo.Client {
	m, _ := mongodb.New(ctx, cfg)
	return m.Client
}

func getCurrentFilename() string {
	_, filename, _, _ := runtime.Caller(0)
	extension := filepath.Ext(filename)
	return strings.TrimSuffix(filepath.Base(filename), extension)
}
