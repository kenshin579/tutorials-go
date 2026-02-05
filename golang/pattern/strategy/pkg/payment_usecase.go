package pkg

import (
	"context"
	"fmt"
)

// PaymentUsecase 는 결제 유스케이스입니다.
// Strategy 를 주입받아 결제를 처리합니다.
type PaymentUsecase struct {
	strategy PaymentStrategy
}

func NewPaymentUsecase(strategy PaymentStrategy) *PaymentUsecase {
	return &PaymentUsecase{strategy: strategy}
}

// SetStrategy 는 런타임에 전략을 변경합니다.
func (u *PaymentUsecase) SetStrategy(strategy PaymentStrategy) {
	u.strategy = strategy
}

// ProcessPayment 는 설정된 전략으로 결제를 처리합니다.
func (u *PaymentUsecase) ProcessPayment(ctx context.Context, amount int) (PaymentResult, error) {
	fmt.Printf("Processing payment with strategy: %s\n", u.strategy.GetName())
	return u.strategy.Pay(ctx, amount)
}

// ProcessRefund 는 설정된 전략으로 환불을 처리합니다.
func (u *PaymentUsecase) ProcessRefund(ctx context.Context, transactionID string) error {
	return u.strategy.Refund(ctx, transactionID)
}
