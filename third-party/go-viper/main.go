package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDBConfig MongoDBConfig `yaml:"mongodb"`
	ServerConfig  ServerConfig  `yaml:"server"`
}

type MongoDBConfig struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

// todo: 잘 안됨 - 나중에 추가하는 걸로 함
func main() {
	//viper.SetConfigName("config") //config 파일 이름 (확장명 제외)
	//viper.SetConfigType("yaml")
	//viper.SetConfigFile("config.yaml")
	//viper.AddConfigPath("./config/")
	viper.SetConfigName("config") //config 파일 이름 (확장명 제외)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./go-viper/config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatalf("%s", err)
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Error reading config file, %s", err)
		}
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	log.Printf("database uri is %s", config.MongoDBConfig.Uri)
	log.Printf("port for this application is %d", config.ServerConfig.Port)

}
