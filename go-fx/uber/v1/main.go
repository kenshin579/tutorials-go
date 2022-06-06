package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "[ACME]", 0)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Handler called")
		io.WriteString(w, "Hello World\n")
	})

	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
