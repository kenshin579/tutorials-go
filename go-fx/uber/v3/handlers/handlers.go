package handlers

import (
	"io"
	"net/http"

	"github.com/kenshin579/tutorials-go/go-fx/uber/v3/loggerfx"
)

func NewHandler(logger loggerfx.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Handler called")
		io.WriteString(w, "Hello World\n")
	})
}
