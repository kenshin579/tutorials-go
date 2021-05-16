package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	fmt.Println("starting scheduler")
	file, err := os.OpenFile("logs-second.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	c := cron.New(
		cron.WithSeconds(),
		cron.WithLogger(cron.VerbosePrintfLogger(log.New(file, "cron: ", log.LstdFlags))))

	//초단위로 잘 되는거 확인함
	c.AddFunc("* * 09-18 * * MON-FRI", func() {
		log.Println("hello world - seconds")
	})
	c.Start()

	time.Sleep(time.Hour * 5)

}
