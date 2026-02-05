package payment

import "context"

// PaymentStrategy 는 결제 처리를 위한 전략 인터페이스입니다.
type PaymentStrategy interface {
	Pay(ctx context.Context, amount int) (PaymentResult, error)
	Refund(ctx context.Context, transactionID string) error
	GetName() string
}

type PaymentResult struct {
	TransactionID string
	Status        string
	Message       string
}
