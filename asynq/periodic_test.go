package asynq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/hibiken/asynq"
	"github.com/kenshin579/tutorials-go/asynq/tasks"
	"github.com/stretchr/testify/assert"
)

func Test_Periodic_Tasks(t *testing.T) {
	redisOpt := asynq.RedisClientOpt{Addr: redisAddr}
	ctx := context.Background()
	redisClient := newRedisClient()
	defer func() {
		redisClient.Close()
	}()

	maxGoroutine := 1

	type schedulerInfo struct {
		scheduler *asynq.Scheduler
		entryID   string
	}

	schedulerInfoMap := make(map[int]schedulerInfo)

	for i := 1; i <= maxGoroutine; i++ {
		go func(i int) {
			scheduler := asynq.NewScheduler(redisOpt, nil)
			// schedulers = append(schedulers, scheduler)

			// You can use cron spec string to specify the schedule.
			loggingTask, err := tasks.NewLoggingTask("title1")
			assert.NoError(t, err)

			// instance 2개: 이렇게 하면 실제로 2초마다 실행되지 않고 4촘마다 실행이 됨
			// entryID, err := scheduler.Register("@every 2s", loggingTask, asynq.TaskID("job1"), asynq.Retention(2*time.Second))
			// instance 2개: trigger도 잘되는 거 확인함
			entryID, err := scheduler.Register("@every 1.5s", loggingTask, asynq.Unique(1500*time.Millisecond), asynq.Retention(5*time.Minute))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("registered an entry: %q\n", entryID)

			schedulerInfoMap[i] = schedulerInfo{scheduler: scheduler, entryID: entryID}

			if err := scheduler.Start(); err != nil {
				log.Fatal(err)
			}
			log.Printf("running scheduler...")

			// 주기 작업은 Register 등록하면 cron library에 의해서 실행이 되고 주기 작업 정보는 redis에 5초마다 쓰여진다
			time.Sleep(1 * time.Second)
			keys := redisClient.Keys(ctx, "asynq:schedulers*")
			assert.NoError(t, err)
			fmt.Printf("keys: %v\n", keys)

		}(i)
	}

	time.Sleep(60 * time.Second)

	for _, schedulerInfo := range schedulerInfoMap {
		assert.NoError(t, schedulerInfo.scheduler.Unregister(schedulerInfo.entryID))
		schedulerInfo.scheduler.Shutdown()
	}

}

var cfgs []*asynq.PeriodicTaskConfig

// Periodic Task manager로 실행을 해도 Scheduler Entries에 추가되는 건 중복으로 들어간다
func Test_Periodic_Task_Manager(t *testing.T) {
	redisOpt := asynq.RedisClientOpt{Addr: redisAddr}

	type schedulerInfo struct {
		scheduler *asynq.Scheduler
		manager   *asynq.PeriodicTaskManager
	}

	maxGoroutine := 2
	schedulerInfoMap := make(map[int]schedulerInfo)

	for i := 1; i <= maxGoroutine; i++ {
		go func(i int) {
			scheduler := asynq.NewScheduler(redisOpt, nil)

			if err := scheduler.Start(); err != nil {
				log.Fatal(err)
			}

			loggingTask1, err := tasks.NewLoggingTask("foo")
			assert.NoError(t, err)

			cfgs = []*asynq.PeriodicTaskConfig{
				{Cronspec: "* * * * *", Task: loggingTask1}, // cron 표현식에서 second이 지원을 하지 않는 듯함
			}

			const syncInterval = 3 * time.Second
			provider := &FakeConfigProvider{cfgs: cfgs}
			manager, err := asynq.NewPeriodicTaskManager(asynq.PeriodicTaskManagerOpts{
				RedisConnOpt:               redisOpt,
				PeriodicTaskConfigProvider: provider,
				SyncInterval:               syncInterval,
			})
			assert.NoError(t, err)

			schedulerInfoMap[i] = schedulerInfo{scheduler: scheduler, manager: manager}

			if err := manager.Start(); err != nil {
				t.Fatalf("Failed to start PeriodicTaskManager: %v", err)
			}

			log.Printf("%d.running task manager", i)
		}(i)
	}

	time.Sleep(60 * time.Second)

	for _, scheduleInfo := range schedulerInfoMap {
		scheduleInfo.manager.Shutdown()
		scheduleInfo.scheduler.Shutdown()
	}

}

type FakeConfigProvider struct {
	mu   sync.Mutex
	cfgs []*asynq.PeriodicTaskConfig
}

func (p *FakeConfigProvider) SetConfigs(cfgs []*asynq.PeriodicTaskConfig) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cfgs = cfgs
}

func (p *FakeConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cfgs, nil
}

func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}
