package pkg

import "context"

// NotificationStrategy는 알림 전송을 위한 전략 인터페이스입니다.
type NotificationStrategy interface {
	Send(ctx context.Context, to string, message string) error
	GetType() NotificationType
}

type NotificationType string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypeSMS   NotificationType = "sms"
	NotificationTypePush  NotificationType = "push"
	NotificationTypeSlack NotificationType = "slack"
)

// NotificationStrategyProvider는 알림 타입에 따라 적절한 Strategy를 제공합니다.
type NotificationStrategyProvider interface {
	Get(ctx context.Context, notificationType NotificationType) (NotificationStrategy, bool)
	GetByUserPreference(ctx context.Context, userID string) (NotificationStrategy, bool)
}

// UserPreferenceStore는 사용자 알림 설정을 조회하는 인터페이스입니다.
type UserPreferenceStore interface {
	GetPreferredNotificationType(ctx context.Context, userID string) (NotificationType, error)
}
