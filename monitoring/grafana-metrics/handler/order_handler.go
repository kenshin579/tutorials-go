package handler

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"grafana-metrics/metrics"
)

// Order 주문 모델
type Order struct {
	ID        string    `json:"id"`
	Product   string    `json:"product"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// 인메모리 저장소
var (
	orders   = make(map[string]*Order)
	ordersMu sync.RWMutex
)

// CreateOrder 주문 생성 핸들러 (랜덤 지연/에러 포함)
func CreateOrder(c echo.Context) error {
	start := time.Now()

	// 랜덤 지연 시뮬레이션 (50~500ms)
	delay := time.Duration(50+rand.Intn(450)) * time.Millisecond
	time.Sleep(delay)

	// 약 10% 확률로 에러 발생
	if rand.Float64() < 0.1 {
		metrics.OrdersCreatedTotal.WithLabelValues("failed").Inc()
		metrics.OrderProcessingDuration.Observe(time.Since(start).Seconds())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "order processing failed",
		})
	}

	order := &Order{
		ID:        uuid.New().String(),
		Product:   "sample-product",
		Amount:    rand.Intn(10000) + 1000,
		Status:    "created",
		CreatedAt: time.Now(),
	}

	ordersMu.Lock()
	orders[order.ID] = order
	ordersMu.Unlock()

	metrics.OrdersCreatedTotal.WithLabelValues("success").Inc()
	metrics.OrderProcessingDuration.Observe(time.Since(start).Seconds())

	return c.JSON(http.StatusCreated, order)
}

// ListOrders 주문 목록 조회 핸들러
func ListOrders(c echo.Context) error {
	// 랜덤 지연 (10~50ms)
	time.Sleep(time.Duration(10+rand.Intn(40)) * time.Millisecond)

	ordersMu.RLock()
	defer ordersMu.RUnlock()

	result := make([]*Order, 0, len(orders))
	for _, o := range orders {
		result = append(result, o)
	}

	return c.JSON(http.StatusOK, result)
}

// GetOrder 주문 상세 조회 핸들러
func GetOrder(c echo.Context) error {
	id := c.Param("id")

	// 랜덤 지연 (5~30ms)
	time.Sleep(time.Duration(5+rand.Intn(25)) * time.Millisecond)

	ordersMu.RLock()
	order, exists := orders[id]
	ordersMu.RUnlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "order not found",
		})
	}

	return c.JSON(http.StatusOK, order)
}
