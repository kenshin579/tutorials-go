package waitgroup

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v9"
	"github.com/kenshin579/tutorials-go/go-concurrency/waitgroup/counter"
	"github.com/kenshin579/tutorials-go/test/localdb"
	"github.com/stretchr/testify/suite"
)

type counterTestSuite struct {
	suite.Suite

	redisClient *redis.Client
	ctx         context.Context
}

func TestCounterTestSuite(t *testing.T) {
	suite.Run(t, new(counterTestSuite))
}

func (suite *counterTestSuite) SetupSuite() {
	fmt.Println("counterTestSuite started")

	redisClient := localdb.NewRedisClient()
	suite.ctx = context.TODO()
	suite.redisClient = redisClient

}

func (suite *counterTestSuite) TearDownTest() {
	suite.NoError(suite.redisClient.FlushAll(suite.ctx).Err())
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

	c.Display() // c의 값 출력
}

// todo: 이게 잘 안됨
func (suite *counterTestSuite) TestCounterRedisLock() {
	// 모든 CPU를 사용하도록 함
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := counter.CounterRedisLock{
		Num:    0,
		Locker: redislock.New(suite.redisClient),
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

	c.Display() // c의 값 출력
}

func (suite *counterTestSuite) TestCounterRedSync() {
	// 모든 CPU를 사용하도록 함
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := counter.CounterRedSync{
		Num:    0,
		Locker: redislock.New(suite.redisClient),
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

	c.Display() // c의 값 출력
}
