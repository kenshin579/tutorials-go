package notification

import (
	"context"
	"fmt"
)

// EmailService 는 이메일 알림 서비스입니다.
type EmailService struct {
	smtpHost string
	smtpPort int
}

func NewEmailService(host string, port int) *EmailService {
	return &EmailService{smtpHost: host, smtpPort: port}
}

func (s *EmailService) Send(ctx context.Context, to string, message string) error {
	fmt.Printf("[Email] Sending to %s via %s:%d\n", to, s.smtpHost, s.smtpPort)
	fmt.Printf("[Email] Message: %s\n", message)
	return nil
}

func (s *EmailService) GetType() NotificationType {
	return NotificationTypeEmail
}
