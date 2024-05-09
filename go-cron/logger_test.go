package go_cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/stretchr/testify/assert"
)

func Test_logger(t *testing.T) {
	logger := gocron.NewLogger(gocron.LogLevelInfo)

	scheduler, _ := gocron.NewScheduler(
		gocron.WithLogger(logger),
	)
	defer func() { _ = scheduler.Shutdown() }()

	job, err := scheduler.NewJob(
		gocron.DurationJob(time.Second),
		gocron.NewTask(func() {
			fmt.Println("running job2")
		}),
	)

	assert.NoError(t, err)

	scheduler.Start()
	nextRun, _ := job.NextRun()
	fmt.Printf("id:%s, name:%s, nextTime:%v\n", job.ID(), job.Name(), nextRun)

	// block until you are ready to shut down
	select {
	case <-time.After(5 * time.Second):
	}

	lastRun, err := job.LastRun()
	fmt.Printf("lastRun:%v, err:%v\n", lastRun, err)

	_ = scheduler.StopJobs()

}
