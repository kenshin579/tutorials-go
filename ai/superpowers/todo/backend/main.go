// Package main is the entrypoint for the superpowers todo learning sample.
// It wires the Echo server with an in-memory todo store and starts listening on :8080.
package main

import (
	"log"

	"github.com/kenshin579/tutorials-go/ai/superpowers/todo/backend/server"
	"github.com/kenshin579/tutorials-go/ai/superpowers/todo/backend/todo"
)

func main() {
	store := todo.NewStore()
	srv := server.New(store)
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
