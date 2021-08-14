package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kenshin579/tutorials-go/go-mongo/adapter/mongodb"
	"github.com/kenshin579/tutorials-go/go-mongo/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/benweissmann/memongo"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/tutorials-go/go-mongo/domain"
)

var (
	store       *mongoStore
	ctx         context.Context
	trainerList = domain.CreateSampleTrainer()
	mongoServer memongo.Server
)

//실제 mongo 서버
func setup() {
	cfg, err := config.New("../../../config/config.yaml")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(cfg)

	mongodb, err := mongodb.New(context.Background(), cfg)
	if err != nil {
		log.Panic(err)
	}
	store = NewMongoStore(mongodb)
	ctx = context.Background()
}

func teardown() {
	store.db.Drop(ctx)
}

func TestMain(m *testing.M) {
	mongoServer, err := memongo.Start("4.0.5")
	clientOpts := options.Client().ApplyURI(mongoServer.URI())
	client, err := mongo.Connect(context.Background(), clientOpts)
	randomDB := client.Database(memongo.RandomDatabase())
	store = NewMongoStore(&mongodb.Mongodb{
		Client: client,
		DB:     randomDB,
	})

	ctx = context.Background()

	if err != nil {
		log.Fatal(err)
	}
	defer mongoServer.Stop()

	os.Exit(m.Run())
}

func TestInsert_FindOne(t *testing.T) {
	err := store.Insert(ctx, trainerList[0])
	assert.NoError(t, err)

	result, err := store.FindOne(ctx, bson.D{{"name", "Ash"}})
	assert.NoError(t, err)
	assert.Equal(t, "Ash", result.Name)

}

func TestInsertMany_FindAll(t *testing.T) {

	err := store.InsertMany(ctx, trainerList)
	assert.NoError(t, err)
	findOptions := options.Find()
	//findOptions.SetLimit(2) //최대 검색 객수 2개로 제한함

	list, err := store.FindAll(ctx, findOptions)
	assert.NoError(t, err)
	assert.Equal(t, len(trainerList), len(list))
}

func TestUpdate(t *testing.T) {

	err := store.Insert(ctx, trainerList[0])
	assert.NoError(t, err)

	result, _ := store.FindOne(ctx, bson.D{{"name", "Ash"}})
	assert.Equal(t, 10, result.Age)

	err = store.Update(ctx,
		bson.D{{"name", "Ash"}},
		bson.D{
			{"$inc", bson.D{
				{"age", 1},
			}},
		})
	assert.NoError(t, err)

	result, err = store.FindOne(ctx, bson.D{{"name", "Ash"}})
	assert.NoError(t, err)
	assert.Equal(t, 11, result.Age)
}

func TestDeleteMany(t *testing.T) {
	err := store.InsertMany(ctx, trainerList)
	assert.NoError(t, err)

	err = store.Delete(ctx, bson.D{{"name", "Ash"}})
	assert.NoError(t, err)
	findOptions := options.Find()
	//findOptions.SetLimit(2) //최대 검색 객수 2개로 제한함

	list, err := store.FindAll(ctx, findOptions)
	assert.Equal(t, len(trainerList)-1, len(list))
}
