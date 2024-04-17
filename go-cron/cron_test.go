package go_cron

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_EventListeners_AfterJonRuns(t *testing.T) {
	scheduler, _ := gocron.NewScheduler()
	defer func() { _ = scheduler.Shutdown() }()

	job, err := scheduler.NewJob(
		gocron.DurationJob(
			time.Second,
		),
		gocron.NewTask(
			func(name string, num int) {
				fmt.Printf("running job. name:%s, num:%d\n", name, num)
			},
			"hello",
			1,
		),
		gocron.WithTags("foo", "bar"),
		gocron.WithEventListeners(
			gocron.AfterJobRuns(
				func(jobID uuid.UUID, jobName string) {
					// do something after the job completes
					fmt.Printf("jobID:%s, jobName:%s\n", jobID, jobName)
				},
			),
		),
	)

	assert.NoError(t, err)
	fmt.Println(job.ID())
	fmt.Println(job.Name())
	fmt.Println(job.Tags())

	scheduler.Start()
	next, _ := job.NextRun()
	fmt.Println(next)

	// Runs the job one time now, without impacting the schedule
	fmt.Println("STARTING JOB NOW")
	_ = job.RunNow()

	// block until you are ready to shut down
	select {
	case <-time.After(5 * time.Second):
	}

	_ = scheduler.StopJobs()
}

func Test_EventListeners_AfterJobRunsWithError(t *testing.T) {
	scheduler, _ := gocron.NewScheduler()
	defer func() { _ = scheduler.Shutdown() }()

	job, err := scheduler.NewJob(
		gocron.DurationJob(
			time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("running job")
			},
		),
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					fmt.Println(jobID, jobName, err)
				},
			),
		),
	)

	assert.NoError(t, err)
	fmt.Println(job.ID())
}

// WithStopTimeout이 잘 동작을 하는지 잘 모르겠음
func Test_WithStopTimeout(t *testing.T) {
	scheduler, _ := gocron.NewScheduler(
		gocron.WithStopTimeout(time.Second * 2),
	)
	defer func() { _ = scheduler.Shutdown() }()

	job, err := scheduler.NewJob(
		gocron.DurationJob(
			time.Second,
		),
		gocron.NewTask(
			func(name string, num int) {
				fmt.Printf("running job. name:%s, num:%d\n", name, num)
				fmt.Println("sleeping for 3 seconds")
				time.Sleep(4 * time.Second)
			},
			"hello",
			1,
		),
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					// do something after the job completes
					fmt.Printf("jobID:%s, jobName:%s, err:%v\n", jobID, jobName, err)
				},
			),
		),
	)

	assert.NoError(t, err)
	next, _ := job.NextRun()
	fmt.Println(next)

	// block until you are ready to shut down
	select {
	case <-time.After(5 * time.Second):
	}

	_ = scheduler.StopJobs()
}

func Test_CronJob(t *testing.T) {
	scheduler, _ := gocron.NewScheduler()
	defer func() { _ = scheduler.Shutdown() }()

	job, err := scheduler.NewJob(
		gocron.CronJob(
			// standard cron tab parsing
			"* * * * * *",
			true,
		),
		gocron.NewTask(
			func() {
				fmt.Println("running job")
			},
		),
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

func Test_OneTimeJob(t *testing.T) {
	scheduler, _ := gocron.NewScheduler()
	defer func() { _ = scheduler.Shutdown() }()

	// run a job once, immediately
	_, _ = scheduler.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartImmediately(),
		),
		gocron.NewTask(
			func() {
				fmt.Println("running task1")
			},
		),
	)
	// run a job once in 3 seconds
	_, _ = scheduler.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(time.Now().Add(3*time.Second)),
		),
		gocron.NewTask(
			func() {
				fmt.Println("running task2")
			},
		),
	)

	scheduler.Start()

	select {
	case <-time.After(5 * time.Second):
	}

	_ = scheduler.StopJobs()
}

func Test_RemoveByTags(t *testing.T) {
	scheduler, _ := gocron.NewScheduler()
	defer func() { _ = scheduler.Shutdown() }()

	_, _ = scheduler.NewJob(
		gocron.DurationJob(
			time.Second,
		),
		gocron.NewTask(
			func() {},
		),
		gocron.WithTags("tag1"),
	)
	_, _ = scheduler.NewJob(
		gocron.DurationJob(
			time.Second,
		),
		gocron.NewTask(
			func() {},
		),
		gocron.WithTags("tag2"),
	)

	assert.Len(t, scheduler.Jobs(), 2)

	scheduler.RemoveByTags("tag1", "tag2")
	assert.Len(t, scheduler.Jobs(), 0)
}

var _ gocron.Locker = (*myLocker)(nil)
var _ gocron.Lock = (*testLock)(nil)

type myLocker struct {
}
type testLock struct {
}

func (m myLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
	return &testLock{}, nil
}

func (t testLock) Unlock(_ context.Context) error {
	return nil
}

func Test_DistributedLocker(t *testing.T) {
	locker := &myLocker{}

	_, _ = gocron.NewScheduler(
		gocron.WithDistributedLocker(locker),
	)
}
