package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	JWTSecret          string
	FrontendURL        string
	ServerPort         string
}

func Load() *Config {
	godotenv.Load() // .env 파일이 없어도 무시

	return &Config{
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:3000/auth/callback"),
		JWTSecret:          getEnv("JWT_SECRET", "dev-only-change-me"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
