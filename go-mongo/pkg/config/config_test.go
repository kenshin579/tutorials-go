package config

import (
	"fmt"
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

	fmt.Println(cfg)

	assert.True(t, strings.ContainsAny("localhost", cfg.MongoConfig.Uri))
	assert.NotEmpty(t, cfg.MongoConfig.Database)
}
