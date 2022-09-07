package waitgroup

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/bsm/redislock"
	redislib_v8 "github.com/go-redis/redis/v8"
	redislib_v9 "github.com/go-redis/redis/v9"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/kenshin579/tutorials-go/go-concurrency/waitgroup/counter"
	"github.com/kenshin579/tutorials-go/test/localdb"
	"github.com/stretchr/testify/suite"
)

type counterTestSuite struct {
	suite.Suite

	redisV8Client *redislib_v8.Client
	redisV9Client *redislib_v9.Client
	ctx           context.Context
}

func TestCounterTestSuite(t *testing.T) {
	suite.Run(t, new(counterTestSuite))
}

func (suite *counterTestSuite) SetupSuite() {
	fmt.Println("counterTestSuite started")

	redisV8Client := localdb.NewRedisV8Client()
	redisV9Client := localdb.NewRedisV9Client()
	suite.ctx = context.TODO()
	suite.redisV8Client = redisV8Client
	suite.redisV9Client = redisV9Client

}

func (suite *counterTestSuite) TearDownTest() {
	suite.NoError(suite.redisV8Client.FlushAll(suite.ctx).Err())
	suite.NoError(suite.redisV9Client.FlushAll(suite.ctx).Err())
}

func (suite *counterTestSuite) TestCounterMutex() {
	// 모든 CPU를 사용하도록 함
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := counter.CounterMutex{Num: 0} // 카운터 생성
	wg := sync.WaitGroup{}            // WaitGroup 생성

	// c.increment()를 실행하는 1000개의 고루틴 실행
	for i := 0; i < 1000; i++ {
		wg.Add(1) // WaitGroup의 고루틴 개수 1 증가
		go func() {
			defer wg.Done() // 고루틴 종료시 Done() 처리
			c.Increment()   // 카운터 값을 1 증가시킴
		}()
	}

	wg.Wait() // 모든 고루틴이 종료될 때 까지 대기

	suite.Equal(int64(1000), c.GetNum())
}

func (suite *counterTestSuite) TestCounterRedisLock() {
	// 모든 CPU를 사용하도록 함
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := counter.CounterRedisLock{
		Num:    0,
		Locker: redislock.New(suite.redisV9Client),
	} // 카운터 생성
	wg := sync.WaitGroup{} // WaitGroup 생성

	// c.increment()를 실행하는 1000개의 고루틴 실행
	for i := 0; i < 1000; i++ {
		wg.Add(1) // WaitGroup의 고루틴 개수 1 증가
		go func() {
			defer wg.Done() // 고루틴 종료시 Done() 처리
			c.Increment()   // 카운터 값을 1 증가시킴
		}()
	}

	wg.Wait() // 모든 고루틴이 종료될 때 까지 대기

	suite.Equal(int64(1000), c.GetNum())
}

func (suite *counterTestSuite) TestCounterRedSync() {
	// 모든 CPU를 사용하도록 함
	runtime.GOMAXPROCS(runtime.NumCPU())

	pool := goredis.NewPool(suite.redisV8Client)
	rs := redsync.New(pool)

	c := counter.CounterRedSync{
		Num:   0,
		Mutex: rs.NewMutex("engine"),
	}

	wg := sync.WaitGroup{} // WaitGroup 생성

	// c.increment()를 실행하는 1000개의 고루틴 실행
	for i := 0; i < 1000; i++ {
		wg.Add(1) // WaitGroup의 고루틴 개수 1 증가
		go func() {
			defer wg.Done() // 고루틴 종료시 Done() 처리
			c.Increment()   // 카운터 값을 1 증가시킴
		}()
	}

	wg.Wait() // 모든 고루틴이 종료될 때 까지 대기

	suite.Equal(int64(1000), c.GetNum())
}
