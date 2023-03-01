package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.HandleFunc("/", index)
	fmt.Println("now serving localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, this is a web page. now:%s", time.Now())
}
