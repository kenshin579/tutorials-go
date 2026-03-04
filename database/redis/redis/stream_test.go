package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newStreamClient(t *testing.T) *redis.Client {
	t.Helper()
	s := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: s.Addr()})
}

func Test_Stream_XAdd_XRead(t *testing.T) {
	rdb := newStreamClient(t)
	ctx := context.Background()

	// XADD - 메시지 추가
	id1, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "mystream",
		Values: map[string]interface{}{
			"user":   "frank",
			"action": "login",
		},
	}).Result()
	assert.NoError(t, err)
	assert.NotEmpty(t, id1)

	id2, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "mystream",
		Values: map[string]interface{}{
			"user":   "angela",
			"action": "purchase",
		},
	}).Result()
	assert.NoError(t, err)
	assert.NotEmpty(t, id2)

	// XREAD - 메시지 읽기 (처음부터)
	streams, err := rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"mystream", "0"},
		Count:   10,
	}).Result()
	assert.NoError(t, err)
	assert.Len(t, streams, 1)
	assert.Len(t, streams[0].Messages, 2)

	// 첫 번째 메시지 확인
	assert.Equal(t, "frank", streams[0].Messages[0].Values["user"])
	assert.Equal(t, "login", streams[0].Messages[0].Values["action"])

	// 두 번째 메시지 확인
	assert.Equal(t, "angela", streams[0].Messages[1].Values["user"])
	assert.Equal(t, "purchase", streams[0].Messages[1].Values["action"])
}

func Test_Stream_XRange(t *testing.T) {
	rdb := newStreamClient(t)
	ctx := context.Background()

	// 메시지 추가
	for i := 0; i < 5; i++ {
		rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: "events",
			Values: map[string]interface{}{
				"event_id": i,
				"type":     "click",
			},
		})
	}

	// XRANGE - 범위 조회 (전체)
	messages, err := rdb.XRange(ctx, "events", "-", "+").Result()
	assert.NoError(t, err)
	assert.Len(t, messages, 5)

	// XREVRANGE - 역순 조회 (최신 2개)
	latest, err := rdb.XRevRangeN(ctx, "events", "+", "-", 2).Result()
	assert.NoError(t, err)
	assert.Len(t, latest, 2)

	// XLEN - 스트림 길이
	length, err := rdb.XLen(ctx, "events").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(5), length)
}

func Test_Stream_ConsumerGroup(t *testing.T) {
	rdb := newStreamClient(t)
	ctx := context.Background()

	stream := "orders"
	group := "order-processors"

	// 메시지 추가
	rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{"order_id": "1001", "item": "laptop"},
	})
	rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{"order_id": "1002", "item": "phone"},
	})
	rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{"order_id": "1003", "item": "tablet"},
	})

	// Consumer Group 생성 (스트림 처음부터 읽기)
	err := rdb.XGroupCreate(ctx, stream, group, "0").Err()
	assert.NoError(t, err)

	// Consumer 1: 메시지 읽기
	result1, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: "consumer-1",
		Streams:  []string{stream, ">"},
		Count:    2,
	}).Result()
	assert.NoError(t, err)
	assert.Len(t, result1[0].Messages, 2)
	assert.Equal(t, "1001", result1[0].Messages[0].Values["order_id"])
	assert.Equal(t, "1002", result1[0].Messages[1].Values["order_id"])

	// Consumer 2: 남은 메시지 읽기
	result2, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: "consumer-2",
		Streams:  []string{stream, ">"},
		Count:    2,
	}).Result()
	assert.NoError(t, err)
	assert.Len(t, result2[0].Messages, 1)
	assert.Equal(t, "1003", result2[0].Messages[0].Values["order_id"])

	// XACK - 메시지 처리 완료 확인
	ackCount, err := rdb.XAck(ctx, stream, group,
		result1[0].Messages[0].ID,
		result1[0].Messages[1].ID,
	).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), ackCount)

	// XPENDING - 미처리 메시지 확인 (consumer-2의 1건)
	pending, err := rdb.XPending(ctx, stream, group).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), pending.Count)
}
