package notification

import (
	"context"
	"fmt"
)

// NotificationUsecase 는 알림 유스케이스입니다.
// Provider 를 통해 적절한 Strategy 를 얻어 알림을 전송합니다.
type NotificationUsecase struct {
	strategyProvider NotificationStrategyProvider
}

func NewNotificationUsecase(provider NotificationStrategyProvider) *NotificationUsecase {
	return &NotificationUsecase{strategyProvider: provider}
}

// SendByType 은 지정된 타입으로 알림을 전송합니다.
func (u *NotificationUsecase) SendByType(ctx context.Context, notificationType NotificationType, to, message string) error {
	strategy, ok := u.strategyProvider.Get(ctx, notificationType)
	if !ok {
		return fmt.Errorf("notification strategy not found for type: %s", notificationType)
	}

	return strategy.Send(ctx, to, message)
}

// SendByUserPreference 는 사용자 설정에 따라 알림을 전송합니다.
func (u *NotificationUsecase) SendByUserPreference(ctx context.Context, userID, message string) error {
	strategy, ok := u.strategyProvider.GetByUserPreference(ctx, userID)
	if !ok {
		return fmt.Errorf("notification strategy not found for user: %s", userID)
	}

	fmt.Printf("Sending notification with strategy: %s\n", strategy.GetType())
	return strategy.Send(ctx, userID, message)
}
