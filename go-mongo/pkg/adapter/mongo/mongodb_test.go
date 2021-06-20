package mongo

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/kenshin579/tutorials-go/go-mongo/pkg/config"
)

var (
	mongodb *Mongodb
)

func setup() {
	cfg, err := config.New("../../../config/config.yaml")
	if err != nil {
		log.Panic("config error: ", err)
	}

	mongodb, err = New(context.Background(), cfg)
	if err != nil {
		log.Panic("error: ", err)
	}
}

func teardown() {
	mongodb.client.Disconnect(context.Background())
}

func TestPing(t *testing.T) {
	setup()
	defer teardown()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := mongodb.client.Ping(ctx, readpref.Primary())
	assert.NoError(t, err)
}
