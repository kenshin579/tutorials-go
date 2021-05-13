package main

import (
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(file, "cron: ", log.LstdFlags))))

	c.AddFunc("* 09-18 * * MON-FRI", func() {
		log.Println("hello world")
	})
	c.Start()

	time.Sleep(time.Hour * 5)

}
