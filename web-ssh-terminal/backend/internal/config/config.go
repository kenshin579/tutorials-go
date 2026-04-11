package config

import (
	"os"

	"web-ssh-terminal/internal/model"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig  `yaml:"server"`
	SSH    SSHConfig     `yaml:"ssh"`
	Robots []model.Robot `yaml:"robots"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type SSHConfig struct {
	PrivateKeyPath string `yaml:"privateKeyPath"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}

	return &cfg, nil
}
