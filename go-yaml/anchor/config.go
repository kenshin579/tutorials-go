package anchor

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	FgColor string `yaml:"fg_color"`
	BgColor string `yaml:"bg_color"`
}

type Colors struct {
	Color1 string `yaml:"color1"`
	Color2 string `yaml:"color2"`
}

func New(configPath string) (Config, error) {
	config := Config{}

	fileBytes, err := os.ReadFile(configPath)
	if err != nil {
		return config, errors.New("fail to read file")
	}

	if err := yaml.Unmarshal(fileBytes, &config); err != nil {
		return config, errors.New("fail to unmarshal")
	}

	return config, nil
}
