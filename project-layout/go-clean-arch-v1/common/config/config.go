package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	ConfigPath = "project-layout/go-clean-arch/config.json"
)

func New() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(ConfigPath)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if v.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}

	return v
}
