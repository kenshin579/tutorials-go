package main

import (
	"fmt"
	"net/http"
	"time"
)

// https://gobyexample.com/context
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	fmt.Println("server: HelloHandler handler started")
	defer fmt.Println("server: HelloHandler handler ended")

	select {
	case <-time.After(1 * time.Second):
		fmt.Fprintf(w, "hello world!")
	case <-ctx.Done():

		err := ctx.Err()
		fmt.Println("server:", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
}

func main() {
	http.HandleFunc("/HelloHandler", HelloHandler)
	http.ListenAndServe(":8090", nil)
}
