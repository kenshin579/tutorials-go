package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newMiniredisClient(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()
	s := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	return rdb, s
}

func Test_Miniredis_Basic(t *testing.T) {
	rdb, _ := newMiniredisClient(t)
	ctx := context.Background()

	// Set/Get
	err := rdb.Set(ctx, "name", "frank", 0).Err()
	assert.NoError(t, err)

	val, err := rdb.Get(ctx, "name").Result()
	assert.NoError(t, err)
	assert.Equal(t, "frank", val)

	// 존재하지 않는 키
	_, err = rdb.Get(ctx, "nonexistent").Result()
	assert.ErrorIs(t, err, redis.Nil)
}

func Test_Miniredis_TTL(t *testing.T) {
	rdb, s := newMiniredisClient(t)
	ctx := context.Background()

	// TTL이 있는 키 설정
	err := rdb.Set(ctx, "session", "token-abc", 10*time.Minute).Err()
	assert.NoError(t, err)

	// 키 존재 확인
	val, err := rdb.Get(ctx, "session").Result()
	assert.NoError(t, err)
	assert.Equal(t, "token-abc", val)

	// FastForward로 시간 경과 시뮬레이션
	s.FastForward(11 * time.Minute)

	// TTL 만료 후 키 사라짐
	_, err = rdb.Get(ctx, "session").Result()
	assert.ErrorIs(t, err, redis.Nil)
}

func Test_Miniredis_SortedSet(t *testing.T) {
	rdb, _ := newMiniredisClient(t)
	ctx := context.Background()

	// 리더보드 데이터 추가
	rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 100, Member: "player1"})
	rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 250, Member: "player2"})
	rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 180, Member: "player3"})

	// 상위 랭킹 조회 (점수 높은 순)
	result, err := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, -1).Result()
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "player2", result[0].Member)
	assert.Equal(t, float64(250), result[0].Score)

	// 특정 멤버 순위 조회 (0-based, 높은 점수 순)
	rank, err := rdb.ZRevRank(ctx, "leaderboard", "player2").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), rank)
}

func Test_Miniredis_PubSub(t *testing.T) {
	rdb, _ := newMiniredisClient(t)
	ctx := context.Background()

	// 구독
	sub := rdb.Subscribe(ctx, "test-channel")
	defer sub.Close()

	_, err := sub.Receive(ctx)
	assert.NoError(t, err)

	// 발행
	rdb.Publish(ctx, "test-channel", "hello miniredis")

	// 수신
	msg, err := sub.ReceiveMessage(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test-channel", msg.Channel)
	assert.Equal(t, "hello miniredis", msg.Payload)
}
