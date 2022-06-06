package server

import (
	"net/http"
)

func RegisterHandlers(handler http.Handler) {
	http.Handle("/", handler)
}

func StartServer() {
	http.ListenAndServe(":8080", nil)
}
