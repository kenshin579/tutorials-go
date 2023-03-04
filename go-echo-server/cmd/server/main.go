package main

import (
	"context"
	"time"

	"github.com/kenshin579/tutorials-go/go-echo-server/cmd/bootstrap"
	"github.com/labstack/gommon/log"
)

// @title Echo API
// @version 1.0
// @description Echo API Server.

// @BasePath /
func main() {
	app := bootstrap.NewApp()

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	<-app.Done()
}
