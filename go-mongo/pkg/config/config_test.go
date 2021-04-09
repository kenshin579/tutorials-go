package config

import (
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestReadConfigFile(t *testing.T) {
	cfg, err := New("config/config.yaml")
	assert.NoError(t, err)

	assert.True(t, strings.ContainsAny("localhost", cfg.MongoDBConfig.Uri))
	assert.NotEmpty(t, cfg.MongoDBConfig.Database)
}
