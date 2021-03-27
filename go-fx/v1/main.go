package main

import (
	"net/http"

	"github.com/kenshin579/tutorials-go/go-fx/server"
)

func main() {
	//1.DI 사용하지 않는 방식
	mux := http.NewServeMux()
	server.New(mux)

	http.ListenAndServe(":8080", mux)
}
