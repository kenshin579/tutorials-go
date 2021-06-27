package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen string `yaml:"listen"`
	Tasks  []struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	} `yaml:"tasks"`
	ServerList []string `yaml:"serverList"`
}

func New(configPath string) (*Config, error) {
	return parseConfigFile(configPath)
}

func parseConfigFile(configPath string) (*Config, error) {
	cfg := &Config{}

	rst, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if yaml.Unmarshal(rst, &cfg); err != nil {
		return nil, fmt.Errorf("failed unmarshal config: %v", err)
	}
	return cfg, nil
}
