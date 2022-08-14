package main

import (
	"net/http"

	"github.com/kenshin579/tutorials-go/go-fx/sumit_agarwal/v1/server"
)

func main() {
	mux := http.NewServeMux()
	server.New(mux)

	http.ListenAndServe(":8080", mux)
}
