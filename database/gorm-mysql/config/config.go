package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MySQL MySQLConfig `yaml:"mysql"`
}

type MySQLConfig struct {
	URL      string `yaml:"url"`
	LogLevel string `yaml:"logLevel"`
}

func ParseFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
