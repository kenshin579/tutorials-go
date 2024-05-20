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
	"github.com/stretchr/testify/assert"
)

const (
	lockKey = "gocron:lock:%s"
)

var (
	ErrFailedToObtainLock  = errors.New("gocron: failed to obtain lock")
	ErrFailedToReleaseLock = errors.New("gocron: failed to release lock")
)

type RedisLocker struct {
	client *redislock.Client
}

func (r *RedisLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
	fmt.Printf("Lock(%s) key:%s\n", time.Now().Format(time.TimeOnly), key)

	redisKey := fmt.Sprintf(lockKey, key)
	lock, err := r.client.Obtain(ctx, redisKey, 6*time.Second, &redislock.Options{ // 여러 스케줄러에서 하나만 실행하기를 원하면 함수의 실행 시간 만큼만 설정하는 게 좋아보임
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 50),
	})
	if err != nil {
		fmt.Printf("failed to obtain lock. key:%s, err:%v\n", key, err)
		return nil, ErrFailedToObtainLock
	}

	return &RedisLock{lock: lock}, nil
}

type RedisLock struct {
	lock *redislock.Lock
}

func (r *RedisLock) Unlock(ctx context.Context) error {
	deadline, ok := ctx.Deadline()
	if ok {
		fmt.Printf("deadline: %s, key:%s\n", deadline.Format(time.TimeOnly), r.lock.Key())
	}

	ttl, err := r.lock.TTL(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get ttl, err:%v", err))
	}
	fmt.Printf("1.ttl:%d, key:%s\n", ttl/time.Second, r.lock.Key())

	fmt.Printf("Unlock(%s), key:%s\n", time.Now().Format(time.TimeOnly), r.lock.Key())
	if err := r.lock.Release(ctx); err != nil { // 프로그램
		fmt.Printf("failed to release lock. key:%s, err:%v\n", r.lock.Key(), err)
		return ErrFailedToReleaseLock
	}
	ttl, err = r.lock.TTL(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get ttl, err:%v", err))
	}
	fmt.Printf("2.ttl:%d, key:%s\n", ttl/time.Second, r.lock.Key())
	return nil
}

func Test_DistributedLockerWithRedis(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	locker := &RedisLocker{
		client: redislock.New(redisClient),
	}

	maxGoroutine := 5
	maxWaitTime := 20 * time.Second
	funRunTime := 5 * time.Second

	resultChan := make(chan int, 10)
	fn1 := func(name string, num int) {
		fmt.Printf("num:%d, START(%s) name:%s\n", num, time.Now().Format(time.TimeOnly), name)
		time.Sleep(funRunTime)
		fmt.Printf("num:%d, END(%s)\n", num, time.Now().Format(time.TimeOnly))
		resultChan <- num
	}

	schedulers := make([]gocron.Scheduler, 0)

	for i := 1; i <= maxGoroutine; i++ {
		go func(i int) {
			scheduler, _ := gocron.NewScheduler(
				gocron.WithDistributedLocker(locker),
			)

			_, err := scheduler.NewJob(
				gocron.DurationJob(
					4*time.Second, // 이 값은 함수 실행 시간 주기로 설정하면 됨 (너무 짧게 설정하면 task를 실행하기 위해 lock을 계속 시도하게 됨)
				),
				gocron.NewTask(fn1, "job1", i),
				// gocron.WithEventListeners(
				// 	gocron.BeforeJobRuns(
				// 		func(jobID uuid.UUID, jobName string) {
				// 			fmt.Printf("%d. beforeJob. jobID:%s, jobName:%s\n", i, jobID, jobName)
				// 		},
				// 	),
				// 	gocron.AfterJobRuns(
				// 		func(jobID uuid.UUID, jobName string) {
				// 			fmt.Printf("%d. jobID:%s, jobName:%s\n", i, jobID, jobName)
				// 		},
				// 	),
				// ),
			)

			assert.NoError(t, err)
			// fmt.Printf("%d. jobName: %v\n", i, job.Name())

			scheduler.Start()
			schedulers = append(schedulers, scheduler)
		}(i)
	}

	select {
	case <-time.After(maxWaitTime):
	}

	for _, s := range schedulers {
		_ = s.StopJobs()
		assert.NoError(t, s.Shutdown())

		// shutdown 할 때 lock을 release 못하는 경우가 발생해서 cleanup을 해줘야 한다
	}

	close(resultChan)

	var results []int
	for r := range resultChan {
		results = append(results, r)
	}
	assert.GreaterOrEqual(t, len(results), int(maxWaitTime.Seconds()/float64(maxGoroutine))-1) // 하나 적게 실행되는 걸 가정함
}
