package cronner

import (
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

type Cronner struct {
	Cron *cron.Cron
}

func New() (*Cronner, error) {
	file, err := os.OpenFile("cronner.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	c := cron.New(
		cron.WithSeconds(),
		cron.WithLogger(cron.VerbosePrintfLogger(log.New(file, "Cron: ", log.LstdFlags))))

	return &Cronner{
		Cron: c,
	}, nil

}
