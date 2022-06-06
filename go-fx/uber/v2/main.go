package main

import (
	"github.com/kenshin579/tutorials-go/go-fx/uber/v2/handlers"
	"github.com/kenshin579/tutorials-go/go-fx/uber/v2/logger"
	"github.com/kenshin579/tutorials-go/go-fx/uber/v2/server"
)

/*
Manual Wiring
*/
func main() {
	logger := logger.NewLogger()
	handler := handlers.NewHandler(logger)
	server.RegisterHandlers(handler)
	server.StartServer()
}
