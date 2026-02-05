package main

import (
	"context"
	"fmt"

	"github.com/kenshin579/tutorials-go/golang/pattern/factory/pkg"
)

func main() {
	ctx := context.Background()

	fmt.Println("========================================")
	fmt.Println("Factory Pattern (Strategy Provider) Example")
	fmt.Println("========================================")
	fmt.Println()

	// 1. 각 Service 생성
	emailService := pkg.NewEmailService("smtp.gmail.com", 587)
	smsService := pkg.NewSMSService("api-key-123", "010-1234-5678")
	pushService := pkg.NewPushService("fcm-token-abc")
	slackService := pkg.NewSlackService("https://hooks.slack.com/xxx")

	// 2. UserPreferenceStore 생성 (사용자 설정 저장소)
	userPrefStore := pkg.NewMockUserPreferenceStore()

	// 3. Strategy Provider (Factory) 생성
	provider := pkg.NewNotificationStrategyProvider(
		emailService,
		smsService,
		pushService,
		slackService,
		userPrefStore,
	)

	// 4. NotificationUsecase 생성
	notificationUsecase := pkg.NewNotificationUsecase(provider)

	// 5. 타입을 직접 지정하여 알림 전송
	fmt.Println("--- 타입 직접 지정 방식 ---")
	_ = notificationUsecase.SendByType(ctx, pkg.NotificationTypeSlack, "#general", "Hello Slack!")
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
