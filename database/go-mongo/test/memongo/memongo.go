package memongo

import (
	"context"
	"runtime"

	"github.com/benweissmann/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMemongoDB() (*mongo.Database, error) {
	opts := &memongo.Options{
		MongoVersion:     "4.2.1",
		ShouldUseReplica: false,
	}

	if runtime.GOARCH == "arm64" {
		if runtime.GOOS == "darwin" {
			// Only set the custom url as workaround for arm64 macs
			opts.DownloadURL = "https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-4.2.1.tgz"
		}
	}

	mongoServer, err := memongo.StartWithOptions(opts)
	if err != nil {
		return nil, err
	}

	clientOpts := options.Client().ApplyURI(mongoServer.URIWithRandomDB())
	mongoClient, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}

	return mongoClient.Database(memongo.RandomDatabase()), nil
}
