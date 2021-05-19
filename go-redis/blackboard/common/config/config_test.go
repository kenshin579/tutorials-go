package config_test

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/kenshin579/tutorials-go/go-redis/blackboard/common/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	fmt.Println(dir)
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestParseFromFile(t *testing.T) {
	cfg, err := config.New("config/config.yaml")
	assert.NoError(t, err)

	assert.NotEmpty(t, cfg.RedisConfig.Addr)
}
