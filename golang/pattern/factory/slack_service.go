package factory

import (
	"context"
	"fmt"
)

// SlackService는 Slack 알림 서비스입니다.
type SlackService struct {
	webhookURL string
}

func NewSlackService(webhookURL string) *SlackService {
	return &SlackService{webhookURL: webhookURL}
}

func (s *SlackService) Send(ctx context.Context, to string, message string) error {
	fmt.Printf("[Slack] Sending to channel %s via webhook\n", to)
	fmt.Printf("[Slack] Message: %s\n", message)
	return nil
}

func (s *SlackService) GetType() NotificationType {
	return NotificationTypeSlack
}
