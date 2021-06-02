package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	fmt.Println(dir)
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestConfigFile(t *testing.T) {
	cfg, err := New("config.yaml")
	assert.NoError(t, err)

	fmt.Println(cfg)

	assert.NotEmpty(t, cfg.RedisConfig)
}
