package account

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
	"github.com/kenshin579/tutorials-go/common/util"
	"github.com/kenshin579/tutorials-go/test"
	"github.com/stretchr/testify/suite"
)

type accountTestSuite struct {
	suite.Suite

	redisV8Client *redislib_v8.Client
	redisV9Client *redislib_v9.Client
	ctx           context.Context
}

func TestCounterTestSuite(t *testing.T) {
	suite.Run(t, new(accountTestSuite))
}

func (suite *accountTestSuite) SetupSuite() {
	fmt.Println("accountTestSuite started")

	redisV8Client := test.NewRedisV8Client()
	redisV9Client := test.NewRedisV9Client()
	suite.ctx = context.TODO()
	suite.redisV8Client = redisV8Client
	suite.redisV9Client = redisV9Client

}

func (suite *accountTestSuite) TearDownTest() {
	suite.NoError(suite.redisV8Client.FlushAll(suite.ctx).Err())
	suite.NoError(suite.redisV9Client.FlushAll(suite.ctx).Err())
}

func (suite *accountTestSuite) TestAccountMutex() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	a := suite.updateBalanceMutix()

	suite.Equal(20000, a.GetBalance("frank"))
	suite.Equal(10000, a.GetBalance("angela"))
}

func (suite *accountTestSuite) updateBalanceMutix() AccountMutex {
	defer util.Timer()("updateBalanceMutix")
	a := AccountMutex{
		customerBalance: map[string]int{"frank": 0, "angela": 0},
	}

	var wg sync.WaitGroup

	doAdd := func(name string, n int) {
		for i := 0; i < n; i++ {
			a.add(name)
		}
		wg.Done()
	}

	wg.Add(3)
	go doAdd("frank", 10000)
	go doAdd("frank", 10000)
	go doAdd("angela", 10000)

	wg.Wait()
	return a
}

func (suite *accountTestSuite) TestAccountRedislock() {
	// 모든 CPU를 사용하도록 함
	runtime.GOMAXPROCS(runtime.NumCPU())

	a := suite.updateBalanceRedislock()

	suite.Equal(20000, a.GetBalance("frank"))
	suite.Equal(10000, a.GetBalance("angela"))
}

func (suite *accountTestSuite) updateBalanceRedislock() AccountRedislock {
	defer util.Timer()("updateBalanceRedislock")

	a := AccountRedislock{
		locker:          redislock.New(suite.redisV9Client),
		customerBalance: map[string]int{"frank": 0, "angela": 0},
	}

	var wg sync.WaitGroup

	doAdd := func(name string, n int) {
		for i := 0; i < n; i++ {
			a.add(name)
		}
		wg.Done()
	}

	wg.Add(3)
	go doAdd("frank", 10000)
	go doAdd("frank", 10000)
	go doAdd("angela", 10000)

	wg.Wait()
	return a
}

func (suite *accountTestSuite) TestAccountRedSync() {
	// 모든 CPU를 사용하도록 함
	runtime.GOMAXPROCS(runtime.NumCPU())

	pool := goredis.NewPool(suite.redisV8Client)
	rs := redsync.New(pool)

	a := suite.updateBalanceRedSync(rs)

	suite.Equal(20000, a.GetBalance("frank"))
	suite.Equal(10000, a.GetBalance("angela"))

}

func (suite *accountTestSuite) updateBalanceRedSync(rs *redsync.Redsync) *AccountRedSync {
	defer util.Timer()("updateBalanceRedSync")

	a := &AccountRedSync{
		Mutex:           rs.NewMutex("account_lock_redsync"),
		CustomerBalance: map[string]int{"frank": 0, "angela": 0},
	}

	var wg sync.WaitGroup

	doAdd := func(name string, n int) {
		for i := 0; i < n; i++ {
			a.add(name)
		}
		wg.Done()
	}

	wg.Add(3)
	go doAdd("frank", 10000)
	go doAdd("frank", 10000)
	go doAdd("angela", 10000)

	wg.Wait()
	return a
}
