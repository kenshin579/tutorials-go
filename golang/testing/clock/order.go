package clock

import "time"

// Order는 생성 시각이 기록되는 주문이다.
type Order struct {
	ID        string
	CreatedAt time.Time
}

// OrderService는 주문을 생성하는 서비스다.
// nowFunc 필드로 현재 시각 함수를 주입받아 테스트에서 시간을 고정할 수 있다.
// 인터페이스 없이 함수 필드 하나로 해결하는 가장 가벼운 주입 패턴이다.
type OrderService struct {
	nowFunc func() time.Time
}

// NewOrderService는 실제 시간(time.Now)을 사용하는 OrderService를 생성한다.
func NewOrderService() *OrderService {
	return &OrderService{nowFunc: time.Now}
}

// CreateOrder는 현재 시각을 생성 시각으로 기록한 주문을 만든다.
func (s *OrderService) CreateOrder(id string) Order {
	return Order{ID: id, CreatedAt: s.nowFunc()}
}
