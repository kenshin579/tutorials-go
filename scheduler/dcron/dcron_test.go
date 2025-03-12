package dcron

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/gochore/dcron"
	"github.com/redis/go-redis/v9"
)

type RedisAtomic struct {
	client *redis.Client
}

func (m *RedisAtomic) SetIfNotExists(ctx context.Context, key, value string) bool {
	ret := m.client.SetNX(ctx, key, value, time.Hour)
	return ret.Err() == nil && ret.Val()
}

func Test_cron(t *testing.T) {
	atomic := &RedisAtomic{
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
	cron := dcron.NewCron(dcron.WithKey("TestCron"), dcron.WithAtomic(atomic))

	job1 := dcron.NewJob("Job1", "*/15 * * * * *", func(ctx context.Context) error {
		if task, ok := dcron.TaskFromContext(ctx); ok {
			log.Println("run:", task.Job.Spec(), task.Key)
		}
		// do something
		return nil
	})
	if err := cron.AddJobs(job1); err != nil {
		log.Fatal(err)
	}

	cron.Start()
	log.Println("cron started")
	time.Sleep(time.Minute)
	<-cron.Stop().Done()
}
