package strategy

import (
	"context"
	"fmt"
)

// KakaoPayService는 카카오페이 결제 서비스입니다.
type KakaoPayService struct {
	userID string
}

func NewKakaoPayService(userID string) *KakaoPayService {
	return &KakaoPayService{userID: userID}
}

func (s *KakaoPayService) Pay(ctx context.Context, amount int) (PaymentResult, error) {
	fmt.Printf("[KakaoPay] Processing payment of %d won for user %s\n", amount, s.userID)
	return PaymentResult{
		TransactionID: "KP-" + generateID(),
		Status:        "SUCCESS",
		Message:       "KakaoPay payment completed",
	}, nil
}

func (s *KakaoPayService) Refund(ctx context.Context, transactionID string) error {
	fmt.Printf("[KakaoPay] Refunding transaction %s\n", transactionID)
	return nil
}

func (s *KakaoPayService) GetName() string {
	return "kakao_pay"
}
