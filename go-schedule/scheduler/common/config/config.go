package config

import (
	"fmt"
	"io/ioutil"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen     string `yaml:"listen"`
	CronConfig []struct {
		Description string                 `yaml:"description"`
		JobType     domain.JobType         `yaml:"jobType"`
		Schedule    string                 `yaml:"schedule"`
		JobRequest  map[string]interface{} `yaml:"jobRequest"`
	} `yaml:"cron"`
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
		return nil, fmt.Errorf("failed unmarshal config: %v", err)
	}
	return cfg, nil
}
