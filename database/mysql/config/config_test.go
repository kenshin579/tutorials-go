package config_test

import (
	"fmt"
	"testing"

	"github.com/kenshin579/tutorials-go/database/mysql/config"
	"github.com/stretchr/testify/assert"
)

func TestParseFromFile(t *testing.T) {
	cfg, err := config.New("config.yaml")
	assert.NoError(t, err)

	fmt.Println(cfg)

	assert.NotEmpty(t, cfg.MysqlConfig.Url)

}
