package main

import (
	"log"
	"net/http"

	"github.com/kenshin579/tutorials-go/go-oom/builder"
	"github.com/kenshin579/tutorials-go/go-oom/validator"
)

func main() {
	go doWork()
	log.Println(http.ListenAndServe("localhost:6060", nil))
}

func doWork() {
	for {
		report := builder.BuildReport()
		validator.ValidateAndSave(report)
	}
}
