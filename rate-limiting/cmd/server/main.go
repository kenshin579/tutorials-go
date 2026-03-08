package main

import (
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kenshin579/tutorials-go/rate-limiting/internal/handler"
	"github.com/kenshin579/tutorials-go/rate-limiting/internal/limiter"
	"github.com/kenshin579/tutorials-go/rate-limiting/internal/middleware"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Default: Token Bucket (10 req capacity, refill 2/sec)
	lim := limiter.NewTokenBucket(rdb, 10, 2.0)

	e := echo.New()
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(middleware.RateLimitMiddleware(lim, middleware.IPKeyFunc))

	handler.RegisterRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s with rate limiting (Token Bucket: capacity=10, refill=2/sec, window=%s)", port, time.Second)
	e.Logger.Fatal(e.Start(":" + port))
}
