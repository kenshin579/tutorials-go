package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/livekit/protocol/auth"
)

const (
	apiKey    = "devkey"
	apiSecret = "secret"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
	}))

	e.GET("/token", handleToken)

	log.Println("Starting token server on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func handleToken(c echo.Context) error {
	roomID := c.QueryParam("roomId")
	userName := c.QueryParam("userName")

	if roomID == "" || userName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "roomId and userName are required")
	}

	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomID,
	}
	at.SetVideoGrant(grant).
		SetIdentity(userName).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
