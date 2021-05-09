package go_schedule

import (
	"testing"

	"github.com/robfig/cron"
)

func Test_cron(t *testing.T) {
	c := cron.New()

	// Every hour on the half hour
	//_ = c.AddFunc("* 30 * * * *", func() {})
	//
	// .. in the range 1-3am, 5-7am
	//_ = c.AddFunc("* 30 1-3, 5-7 * * *", func() {})

	// Every Monday and Thursday
	//_ = c.AddFunc("* * * * * MON,THU", func() {})

	// Every hour, starting an hour from now
	//_ = c.AddFunc("@hourly", func() {})

	// Every hour thirty, starting an hour thirty from now
	//_ = c.AddFunc("@every 1h30m", func() {})
	//c.AddJob()

	c.Start()

}
