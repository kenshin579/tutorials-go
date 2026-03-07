package handler

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"grafana-tracing/metrics"
)

var tracer = otel.Tracer("order-service")

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

// CreateOrder 주문 생성 핸들러 (커스텀 Span 포함)
func CreateOrder(c echo.Context) error {
	ctx := c.Request().Context()
	start := time.Now()

	// 주문 유효성 검증 Span
	order, err := validateOrder(ctx)
	if err != nil {
		metrics.OrdersCreatedTotal.WithLabelValues("failed").Inc()
		metrics.OrderProcessingDuration.Observe(time.Since(start).Seconds())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// 결제 처리 Span
	if err := processPayment(ctx, order); err != nil {
		metrics.OrdersCreatedTotal.WithLabelValues("failed").Inc()
		metrics.OrderProcessingDuration.Observe(time.Since(start).Seconds())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "payment failed"})
	}

	// DB 저장 Span
	if err := saveOrder(ctx, order); err != nil {
		metrics.OrdersCreatedTotal.WithLabelValues("failed").Inc()
		metrics.OrderProcessingDuration.Observe(time.Since(start).Seconds())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "save failed"})
	}

	metrics.OrdersCreatedTotal.WithLabelValues("success").Inc()
	metrics.OrderProcessingDuration.Observe(time.Since(start).Seconds())

	return c.JSON(http.StatusCreated, order)
}

func validateOrder(ctx context.Context) (*Order, error) {
	ctx, span := tracer.Start(ctx, "validate-order")
	defer span.End()

	// 유효성 검증 시뮬레이션 (1~5ms)
	time.Sleep(time.Duration(1+rand.Intn(5)) * time.Millisecond)

	order := &Order{
		ID:        uuid.New().String(),
		Product:   "sample-product",
		Amount:    rand.Intn(10000) + 1000,
		Status:    "created",
		CreatedAt: time.Now(),
	}

	span.SetAttributes(
		attribute.String("order.id", order.ID),
		attribute.Int("order.amount", order.Amount),
		attribute.String("order.product", order.Product),
	)

	return order, nil
}

func processPayment(ctx context.Context, order *Order) error {
	ctx, span := tracer.Start(ctx, "process-payment")
	defer span.End()

	// 결제 처리 시뮬레이션 (50~200ms)
	delay := time.Duration(50+rand.Intn(150)) * time.Millisecond
	time.Sleep(delay)

	span.SetAttributes(
		attribute.String("payment.order_id", order.ID),
		attribute.Int("payment.amount", order.Amount),
	)

	// 약 5% 확률로 결제 실패
	if rand.Float64() < 0.05 {
		span.SetStatus(codes.Error, "payment processing failed")
		span.RecordError(fmt.Errorf("payment gateway timeout"))
		span.AddEvent("payment-failed", trace.WithAttributes(
			attribute.String("reason", "gateway_timeout"),
		))
		return fmt.Errorf("payment failed")
	}

	span.AddEvent("payment-completed")
	return nil
}

func saveOrder(ctx context.Context, order *Order) error {
	_, span := tracer.Start(ctx, "save-order")
	defer span.End()

	// DB 저장 시뮬레이션 (5~20ms)
	time.Sleep(time.Duration(5+rand.Intn(15)) * time.Millisecond)

	span.SetAttributes(
		attribute.String("db.system", "memory"),
		attribute.String("db.operation", "INSERT"),
	)

	ordersMu.Lock()
	orders[order.ID] = order
	ordersMu.Unlock()

	return nil
}

// ListOrders 주문 목록 조회 핸들러
func ListOrders(c echo.Context) error {
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
