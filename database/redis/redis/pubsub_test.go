package redis

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newPubSubClient(t *testing.T) *redis.Client {
	t.Helper()
	s := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: s.Addr()})
}

func Test_PubSub_Basic(t *testing.T) {
	rdb := newPubSubClient(t)
	ctx := context.Background()

	// 채널 구독
	sub := rdb.Subscribe(ctx, "notifications")
	defer sub.Close()

	// 구독 확인 대기
	_, err := sub.Receive(ctx)
	assert.NoError(t, err)

	// 메시지 발행
	err = rdb.Publish(ctx, "notifications", "hello redis pubsub").Err()
	assert.NoError(t, err)

	// 메시지 수신
	msg, err := sub.ReceiveMessage(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "notifications", msg.Channel)
	assert.Equal(t, "hello redis pubsub", msg.Payload)
}

func Test_PubSub_PatternSubscribe(t *testing.T) {
	rdb := newPubSubClient(t)
	ctx := context.Background()

	// 패턴 구독: news.* 채널 전체 매칭
	sub := rdb.PSubscribe(ctx, "news.*")
	defer sub.Close()

	_, err := sub.Receive(ctx)
	assert.NoError(t, err)

	// 다른 채널로 발행
	rdb.Publish(ctx, "news.tech", "Go 1.22 released")
	rdb.Publish(ctx, "news.sports", "World Cup 2026")

	// 두 메시지 모두 수신
	msg1, err := sub.ReceiveMessage(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "news.tech", msg1.Channel)
	assert.Equal(t, "news.*", msg1.Pattern)
	assert.Equal(t, "Go 1.22 released", msg1.Payload)

	msg2, err := sub.ReceiveMessage(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "news.sports", msg2.Channel)
	assert.Equal(t, "World Cup 2026", msg2.Payload)
}

func Test_PubSub_MultiChannel(t *testing.T) {
	rdb := newPubSubClient(t)
	ctx := context.Background()

	// 다중 채널 동시 구독
	sub := rdb.Subscribe(ctx, "channel-1", "channel-2", "channel-3")
	defer sub.Close()

	_, err := sub.Receive(ctx)
	assert.NoError(t, err)

	// 각 채널에 메시지 발행
	messages := map[string]string{
		"channel-1": "msg-1",
		"channel-2": "msg-2",
		"channel-3": "msg-3",
	}

	for ch, msg := range messages {
		err := rdb.Publish(ctx, ch, msg).Err()
		assert.NoError(t, err)
	}

	// Channel()을 사용한 수신 (goroutine 패턴)
	ch := sub.Channel()
	received := make(map[string]string)
	var mu sync.Mutex

	done := make(chan struct{})
	go func() {
		for msg := range ch {
			mu.Lock()
			received[msg.Channel] = msg.Payload
			if len(received) == 3 {
				mu.Unlock()
				close(done)
				return
			}
			mu.Unlock()
		}
	}()

	select {
	case <-done:
		// 성공
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for messages")
	}

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, "msg-1", received["channel-1"])
	assert.Equal(t, "msg-2", received["channel-2"])
	assert.Equal(t, "msg-3", received["channel-3"])
}
