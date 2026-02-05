package payment

import (
	"context"
	"fmt"
)

// CreditCardService 는 신용카드 결제 서비스입니다.
type CreditCardService struct {
	cardNumber string
	cvv        string
}

func NewCreditCardService(cardNumber, cvv string) *CreditCardService {
	return &CreditCardService{
		cardNumber: cardNumber,
		cvv:        cvv,
	}
}

func (s *CreditCardService) Pay(ctx context.Context, amount int) (PaymentResult, error) {
	fmt.Printf("[CreditCard] Processing payment of %d won with card %s\n", amount, s.maskCardNumber())
	return PaymentResult{
		TransactionID: "CC-" + generateID(),
		Status:        "SUCCESS",
		Message:       "Credit card payment completed",
	}, nil
}

func (s *CreditCardService) Refund(ctx context.Context, transactionID string) error {
	fmt.Printf("[CreditCard] Refunding transaction %s\n", transactionID)
	return nil
}

func (s *CreditCardService) GetName() string {
	return "credit_card"
}

func (s *CreditCardService) maskCardNumber() string {
	if len(s.cardNumber) < 4 {
		return "****"
	}
	return "****-****-****-" + s.cardNumber[len(s.cardNumber)-4:]
}
