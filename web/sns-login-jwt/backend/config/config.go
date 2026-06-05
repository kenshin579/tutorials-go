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
		GoogleClientID:     getEnvAny([]string{"WEB_SNS_LOGIN_GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_ID"}, ""),
		GoogleClientSecret: getEnvAny([]string{"WEB_SNS_LOGIN_GOOGLE_CLIENT_SECRET", "GOOGLE_CLIENT_SECRET"}, ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:3000/auth/jwt/callback"),
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

// getEnvAny는 keys를 순서대로 조회하여 처음으로 값이 있는 환경 변수를 반환한다.
func getEnvAny(keys []string, defaultVal string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}
	return defaultVal
}
