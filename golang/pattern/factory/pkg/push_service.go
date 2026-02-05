package pkg

import (
	"context"
	"fmt"
)

// PushService는 푸시 알림 서비스입니다.
type PushService struct {
	fcmToken string
}

func NewPushService(fcmToken string) *PushService {
	return &PushService{fcmToken: fcmToken}
}

func (s *PushService) Send(ctx context.Context, to string, message string) error {
	fmt.Printf("[Push] Sending to device %s\n", to)
	fmt.Printf("[Push] Message: %s\n", message)
	return nil
}

func (s *PushService) GetType() NotificationType {
	return NotificationTypePush
}
