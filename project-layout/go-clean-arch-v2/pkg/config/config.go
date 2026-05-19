package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug    bool     `yaml:"debug"`
	Server   Server   `yaml:"server"`
	Context  Context  `yaml:"context"`
	Database Database `yaml:"database"`
}

type Server struct {
	Address string `yaml:"address"`
}

type Context struct {
	Timeout int `yaml:"timeout"` // seconds
}

type Database struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
}

func New() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
