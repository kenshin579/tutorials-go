package patterns

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Broker - 간단한 Pub/Sub 브로커
type Broker[T any] struct {
	mu          sync.RWMutex      // RWMutex: 다중 Publish는 RLock으로 동시 진행, 구조 변경 시만 Lock
	subscribers map[string]chan T // id → 구독자 channel 맵
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		subscribers: make(map[string]chan T),
	}
}

func (b *Broker[T]) Subscribe(id string, bufSize int) <-chan T {
	b.mu.Lock() // 맵 변경(쓰기)이므로 Lock 필요
	defer b.mu.Unlock()
	ch := make(chan T, bufSize) // bufSize: 느린 구독자가 메시지를 버퍼링할 양
	b.subscribers[id] = ch
	return ch // 수신 전용 반환: 구독자는 받기만, 직접 송신/close 금지
}

func (b *Broker[T]) Unsubscribe(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if ch, ok := b.subscribers[id]; ok {
		close(ch)                 // 구독자에게 종료 신호 (수신 측 range 종료)
		delete(b.subscribers, id) // 맵에서 제거
	}
}

func (b *Broker[T]) Publish(msg T) {
	b.mu.RLock() // 읽기 잠금: 여러 publisher가 동시 발행 가능
	defer b.mu.RUnlock()
	for _, ch := range b.subscribers {
		select {
		case ch <- msg: // 정상 발행
		default:        // 구독자 channel 가득 시 메시지 드롭 → 느린 구독자가 전체를 멈추지 않게 함
		}
	}
}

// TestPubSub - Pub/Sub 패턴 테스트
func TestPubSub(t *testing.T) {
	broker := NewBroker[string]()

	sub1 := broker.Subscribe("sub1", 10)
	sub2 := broker.Subscribe("sub2", 10)

	broker.Publish("hello")
	broker.Publish("world")

	assert.Equal(t, "hello", <-sub1)
	assert.Equal(t, "world", <-sub1)
	assert.Equal(t, "hello", <-sub2)
	assert.Equal(t, "world", <-sub2)

	broker.Unsubscribe("sub1")
	broker.Publish("after unsub")

	// sub2만 받음
	assert.Equal(t, "after unsub", <-sub2)

	// sub1은 닫힘
	_, ok := <-sub1
	assert.False(t, ok)

	broker.Unsubscribe("sub2")
}

// TestPubSubConcurrent - Pub/Sub concurrent 사용
func TestPubSubConcurrent(t *testing.T) {
	broker := NewBroker[int]()
	var wg sync.WaitGroup

	// 3개의 subscriber
	subs := make([]<-chan int, 3)
	for i := range 3 {
		subs[i] = broker.Subscribe(string(rune('a'+i)), 100)
	}

	// publisher
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 10 {
			broker.Publish(i)
			time.Sleep(5 * time.Millisecond)
		}
	}()

	wg.Wait()
	time.Sleep(20 * time.Millisecond)

	for i, sub := range subs {
		count := len(sub)
		t.Logf("subscriber %d received %d messages", i, count)
		assert.Greater(t, count, 0)
	}
}
