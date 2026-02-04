package pattern

import (
	"context"
	"fmt"
	"testing"

	"github.com/kenshin579/tutorials-go/golang/pattern/factory"
	"github.com/kenshin579/tutorials-go/golang/pattern/strategy"
)

// =============================================================================
// Strategy Pattern 예제
// =============================================================================

func TestStrategyPattern(t *testing.T) {
	ctx := context.Background()

	fmt.Println("========================================")
	fmt.Println("Strategy Pattern Example")
	fmt.Println("========================================")
	fmt.Println()

	// 1. 신용카드 결제 서비스로 시작
	creditCardService := strategy.NewCreditCardService("1234567890123456", "123")
	paymentUsecase := strategy.NewPaymentUsecase(creditCardService)

	fmt.Println("--- 신용카드로 결제 ---")
	result, _ := paymentUsecase.ProcessPayment(ctx, 50000)
	fmt.Printf("Result: %+v\n\n", result)

	// 2. 런타임에 카카오페이 서비스로 전략 변경
	kakaoPayService := strategy.NewKakaoPayService("user123")
	paymentUsecase.SetStrategy(kakaoPayService)

	fmt.Println("--- 카카오페이로 결제 ---")
	result, _ = paymentUsecase.ProcessPayment(ctx, 30000)
	fmt.Printf("Result: %+v\n\n", result)

	// 3. 네이버페이 서비스로 전략 변경
	naverPayService := strategy.NewNaverPayService("user456")
	paymentUsecase.SetStrategy(naverPayService)

	fmt.Println("--- 네이버페이로 결제 ---")
	result, _ = paymentUsecase.ProcessPayment(ctx, 25000)
	fmt.Printf("Result: %+v\n\n", result)
}

// =============================================================================
// Factory Pattern (Strategy Provider) 예제
// =============================================================================

func TestFactoryPattern(t *testing.T) {
	ctx := context.Background()

	fmt.Println("========================================")
	fmt.Println("Factory Pattern (Strategy Provider) Example")
	fmt.Println("========================================")
	fmt.Println()

	// 1. 각 Service 생성
	emailService := factory.NewEmailService("smtp.gmail.com", 587)
	smsService := factory.NewSMSService("api-key-123", "010-1234-5678")
	pushService := factory.NewPushService("fcm-token-abc")
	slackService := factory.NewSlackService("https://hooks.slack.com/xxx")

	// 2. UserPreferenceStore 생성 (사용자 설정 저장소)
	userPrefStore := factory.NewMockUserPreferenceStore()

	// 3. Strategy Provider (Factory) 생성
	provider := factory.NewNotificationStrategyProvider(
		emailService,
		smsService,
		pushService,
		slackService,
		userPrefStore,
	)

	// 4. NotificationUsecase 생성
	notificationUsecase := factory.NewNotificationUsecase(provider)

	// 5. 타입을 직접 지정하여 알림 전송
	fmt.Println("--- 타입 직접 지정 방식 ---")
	_ = notificationUsecase.SendByType(ctx, factory.NotificationTypeSlack, "#general", "Hello Slack!")
	fmt.Println()

	// 6. 사용자별 설정에 따라 자동으로 Strategy 선택
	fmt.Println("--- 사용자 설정 기반 방식 (Factory Pattern 핵심) ---")
	fmt.Println()

	users := []string{"user001", "user002", "user003", "user004"}
	for _, userID := range users {
		fmt.Printf("Sending notification to %s:\n", userID)
		err := notificationUsecase.SendByUserPreference(ctx, userID, "Your order has been shipped!")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println()
	}
}

