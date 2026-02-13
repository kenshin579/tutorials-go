package main

import (
	"log"

	"github.com/kenshin579/tutorials-go/webrtc/multi-users-sfu-pion/backend/handler"
	"github.com/kenshin579/tutorials-go/webrtc/multi-users-sfu-pion/backend/room"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
	}))

	rm := room.NewManager()
	sh := handler.NewSignaling(rm)

	e.GET("/ws", sh.HandleWebSocket)

	log.Println("Starting SFU server on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
