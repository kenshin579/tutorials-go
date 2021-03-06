package config_test

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/config"
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

func TestParseConfigFile(t *testing.T) {
	cfg, err := config.New("config/config.yaml")
	assert.NoError(t, err)
	fmt.Println(cfg)

	assert.NotEmpty(t, cfg.CronConfig)
}
