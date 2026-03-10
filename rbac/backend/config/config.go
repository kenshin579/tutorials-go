package config

import (
	"fmt"
	"os"
)

type Config struct {
	DB  DBConfig
	JWT JWTConfig
	App AppConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

type JWTConfig struct {
	Secret string
}

type AppConfig struct {
	Port string
}

func Load() *Config {
	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3307"),
			User:     getEnv("DB_USER", "rbac_user"),
			Password: getEnv("DB_PASSWORD", "rbac_pass"),
			Name:     getEnv("DB_NAME", "rbac_db"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "rbac-sample-secret-key-change-in-production"),
		},
		App: AppConfig{
			Port: getEnv("APP_PORT", "8081"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
