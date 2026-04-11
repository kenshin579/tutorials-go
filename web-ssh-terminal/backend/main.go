package main

import (
	"fmt"
	"log"

	"web-ssh-terminal/internal/config"
	"web-ssh-terminal/internal/handler"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:5174"},
	}))

	robotHandler := handler.NewRobotHandler(cfg)
	terminalHandler := handler.NewTerminalHandler(cfg)

	e.GET("/api/robots", robotHandler.ListRobots)
	e.GET("/ws/terminal", terminalHandler.HandleTerminal)

	// Production: serve React build
	e.Static("/", "../frontend/dist")

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	e.Logger.Fatal(e.Start(addr))
}
