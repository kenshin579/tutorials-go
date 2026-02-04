package factory

import (
	"context"
	"fmt"
)

// notificationStrategyProvider는 NotificationStrategyProvider의 구현체입니다.
type notificationStrategyProvider struct {
	strategies          map[NotificationType]NotificationStrategy
	userPreferenceStore UserPreferenceStore
}

// NewNotificationStrategyProvider는 Provider를 생성합니다.
func NewNotificationStrategyProvider(
	emailService *EmailService,
	smsService *SMSService,
	pushService *PushService,
	slackService *SlackService,
	userPreferenceStore UserPreferenceStore,
) NotificationStrategyProvider {
	strategies := map[NotificationType]NotificationStrategy{
		NotificationTypeEmail: emailService,
		NotificationTypeSMS:   smsService,
		NotificationTypePush:  pushService,
		NotificationTypeSlack: slackService,
	}

	return &notificationStrategyProvider{
		strategies:          strategies,
		userPreferenceStore: userPreferenceStore,
	}
}

// Get은 알림 타입에 해당하는 Strategy를 반환합니다.
func (p *notificationStrategyProvider) Get(ctx context.Context, notificationType NotificationType) (NotificationStrategy, bool) {
	strategy, ok := p.strategies[notificationType]
	return strategy, ok
}

// GetByUserPreference는 사용자 설정에 따른 Strategy를 반환합니다.
func (p *notificationStrategyProvider) GetByUserPreference(ctx context.Context, userID string) (NotificationStrategy, bool) {
	preferredType, err := p.userPreferenceStore.GetPreferredNotificationType(ctx, userID)
	if err != nil {
		fmt.Printf("Failed to get user preference: %v\n", err)
		return nil, false
	}

	strategy, ok := p.strategies[preferredType]
	return strategy, ok
}
