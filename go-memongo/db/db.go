package db

import (
	"context"
	"runtime"

	"github.com/labstack/gommon/log"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewInMemoryMongoDB() *mongo.Database {
	option := getOption()
	mongoServer, err := memongo.StartWithOptions(option)
	if err != nil {
		log.Fatalf(err.Error())
	}

	clientOpts := options.Client().ApplyURI(mongoServer.URI())
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return client.Database(memongo.RandomDatabase())
}

func getOption() *memongo.Options {
	opts := &memongo.Options{
		MongoVersion:     "4.2.1",
		ShouldUseReplica: true,
	}
	return opts
}

// getM1Option 옵션으로 실행할 수 있지만, background에 mongodb를 띄우고
func getM1Option() *memongo.Options {
	opts := &memongo.Options{
		ShouldUseReplica: true,
		MongoVersion:     "4.2.1",
		LogLevel:         2,
	}

	if runtime.GOARCH == "arm64" {
		if runtime.GOOS == "darwin" {
			// Only set the custom url as workaround for arm64 macs
			opts.DownloadURL = "https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-4.2.1.tgz"
		}
	}
	return opts
}
