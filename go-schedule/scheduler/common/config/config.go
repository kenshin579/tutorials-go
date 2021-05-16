package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen       string    `yaml:"listen"`
	ServerConfig ApiConfig `yaml:"server"`
}

type ApiConfig struct {
	Addr string `yaml:"addr"`
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

	if yaml.Unmarshal(rst, &cfg); err != nil {
		return nil, fmt.Errorf("failed to Make arbiter config: %v", err)
	}
	return cfg, nil
}
