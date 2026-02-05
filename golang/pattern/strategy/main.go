package main

import (
	"context"
	"fmt"

	"github.com/kenshin579/tutorials-go/golang/pattern/strategy/payment"
)

func main() {
	ctx := context.Background()

	fmt.Println("========================================")
	fmt.Println("Strategy Pattern Example")
	fmt.Println("========================================")
	fmt.Println()

	// 1. 신용카드 결제 서비스로 시작
	creditCardService := payment.NewCreditCardService("1234567890123456", "123")
	paymentUsecase := payment.NewPaymentUsecase(creditCardService)

	fmt.Println("--- 신용카드로 결제 ---")
	result, _ := paymentUsecase.ProcessPayment(ctx, 50000)
	fmt.Printf("Result: %+v\n\n", result)

	// 2. 런타임에 카카오페이 서비스로 전략 변경
	kakaoPayService := payment.NewKakaoPayService("user123")
	paymentUsecase.SetStrategy(kakaoPayService)

	fmt.Println("--- 카카오페이로 결제 ---")
	result, _ = paymentUsecase.ProcessPayment(ctx, 30000)
	fmt.Printf("Result: %+v\n\n", result)

	// 3. 네이버페이 서비스로 전략 변경
	naverPayService := payment.NewNaverPayService("user456")
	paymentUsecase.SetStrategy(naverPayService)

	fmt.Println("--- 네이버페이로 결제 ---")
	result, _ = paymentUsecase.ProcessPayment(ctx, 25000)
	fmt.Printf("Result: %+v\n\n", result)
}
