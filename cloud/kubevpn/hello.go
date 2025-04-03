package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, "Hello world!")
		fmt.Printf(">>Received request: %s %s from %s\n", request.Method, request.RequestURI, request.RemoteAddr)
	})
	_ = http.ListenAndServe(":9080", nil)
}
