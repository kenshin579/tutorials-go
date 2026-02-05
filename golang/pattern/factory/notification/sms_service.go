package notification

import (
	"context"
	"fmt"
)

// SMSService 는 SMS 알림 서비스입니다.
type SMSService struct {
	apiKey    string
	senderNum string
}

func NewSMSService(apiKey, senderNum string) *SMSService {
	return &SMSService{apiKey: apiKey, senderNum: senderNum}
}

func (s *SMSService) Send(ctx context.Context, to string, message string) error {
	fmt.Printf("[SMS] Sending to %s from %s\n", to, s.senderNum)
	fmt.Printf("[SMS] Message: %s\n", message)
	return nil
}

func (s *SMSService) GetType() NotificationType {
	return NotificationTypeSMS
}
