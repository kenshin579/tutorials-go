package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MysqlConfig MysqlConfig `yaml:"mysql"`
}

type MysqlConfig struct {
	Url      string `yaml:"url"`
	LogLevel string `yaml:"logLevel"`
}

func New(configPath string) (*Config, error) {
	return parseFromFile(configPath)
}

func parseFromFile(configPath string) (*Config, error) {
	cfg := &Config{}

	rst, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(rst, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}
	return cfg, nil
}
