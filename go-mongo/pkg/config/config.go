package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MongoDBConfig MongoDBConfig `yaml:"mongodb"`
}

type MongoDBConfig struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func New(configPath string) (*Config, error) {
	return readConfigFile(configPath)
}

func readConfigFile(configPath string) (*Config, error) {
	cfg := &Config{}

	rst, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if yaml.Unmarshal(rst, &cfg); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	return cfg, nil
}
