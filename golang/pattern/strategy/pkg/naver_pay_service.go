package pkg

import (
	"context"
	"fmt"
)

// NaverPayService 는 네이버페이 결제 서비스입니다.
type NaverPayService struct {
	userID string
}

func NewNaverPayService(userID string) *NaverPayService {
	return &NaverPayService{userID: userID}
}

func (s *NaverPayService) Pay(ctx context.Context, amount int) (PaymentResult, error) {
	fmt.Printf("[NaverPay] Processing payment of %d won for user %s\n", amount, s.userID)
	return PaymentResult{
		TransactionID: "NP-" + generateID(),
		Status:        "SUCCESS",
		Message:       "NaverPay payment completed",
	}, nil
}

func (s *NaverPayService) Refund(ctx context.Context, transactionID string) error {
	fmt.Printf("[NaverPay] Refunding transaction %s\n", transactionID)
	return nil
}

func (s *NaverPayService) GetName() string {
	return "naver_pay"
}
