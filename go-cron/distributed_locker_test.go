package go_cron

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const lockKey = "gocron:lock:%s"

var (
	ErrFailedToObtainLock  = errors.New("gocron: failed to obtain lock")
	ErrFailedToReleaseLock = errors.New("gocron: failed to release lock")
)

type RedisLocker struct {
	client *redislock.Client
}

func (r *RedisLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
	fmt.Printf("[FRANK] Lock. key:%s\n", key)

	redisKey := fmt.Sprintf(lockKey, key)
	lock, err := r.client.Obtain(ctx, redisKey, 100*time.Millisecond, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 50),
	})
	if err != nil {
		return nil, ErrFailedToObtainLock
	}

	return &RedisLock{lock: lock}, nil
}

type RedisLock struct {
	lock *redislock.Lock
}

func (r *RedisLock) Unlock(ctx context.Context) error {
	fmt.Printf("[FRANK] Unlock\n")
	if err := r.lock.Release(ctx); err != nil {
		return ErrFailedToReleaseLock
	}
	return nil
}

func Test_DistributedLockerWithRedis(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	locker := &RedisLocker{
		client: redislock.New(redisClient),
	}

	resultChan := make(chan int, 10)
	fn := func(name string, num int) {
		fmt.Printf("running job. name:%s, num:%d\n", name, num)
		// time.Sleep(2 * time.Second)
		resultChan <- num
	}

	schedulers := make([]gocron.Scheduler, 0)
	for i := 1; i <= 3; i++ {
		go func(i int) {
			scheduler, _ := gocron.NewScheduler(
				gocron.WithDistributedLocker(locker),
			)

			job, err := scheduler.NewJob(
				gocron.DurationJob(
					time.Second,
				),
				gocron.NewTask(fn, "job1", i),
				gocron.WithEventListeners(
					gocron.BeforeJobRuns(
						func(jobID uuid.UUID, jobName string) {
							fmt.Printf("%d. beforeJob. jobID:%s, jobName:%s\n", i, jobID, jobName)
						},
					),
					gocron.AfterJobRuns(
						func(jobID uuid.UUID, jobName string) {
							fmt.Printf("%d. jobID:%s, jobName:%s\n", i, jobID, jobName)
						},
					),
				),
			)

			assert.NoError(t, err)
			fmt.Printf("%d.job: %v\n", i, job.Name())

			scheduler.Start()
			schedulers = append(schedulers, scheduler)
		}(i)
	}

	select {
	case <-time.After(4 * time.Second):
	}

	for _, s := range schedulers {
		_ = s.StopJobs()
		assert.NoError(t, s.Shutdown())
	}

	close(resultChan)

	var results []int
	for r := range resultChan {
		results = append(results, r)
	}
	assert.Len(t, results, 3)
}
