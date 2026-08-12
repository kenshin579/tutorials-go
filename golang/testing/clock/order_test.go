package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_OrderService_고정_시간_주입(t *testing.T) {
	fixed := time.Date(2026, 8, 12, 9, 30, 0, 0, time.UTC)
	svc := &OrderService{nowFunc: func() time.Time { return fixed }}

	order := svc.CreateOrder("order-1")

	assert.Equal(t, fixed, order.CreatedAt)
}

func Test_OrderService_주입_없이_WithinDuration_검증(t *testing.T) {
	svc := NewOrderService()

	order := svc.CreateOrder("order-2")

	// 시간을 주입할 수 없을 때의 차선책: 정확한 일치 대신 오차 허용 비교
	assert.WithinDuration(t, time.Now(), order.CreatedAt, time.Second)
}
