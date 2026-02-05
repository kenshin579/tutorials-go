package notification

import (
	"context"
	"errors"
)

// mockUserPreferenceStore 는 테스트/데모용 UserPreferenceStore 구현체입니다.
type mockUserPreferenceStore struct {
	preferences map[string]NotificationType
}

func NewMockUserPreferenceStore() UserPreferenceStore {
	return &mockUserPreferenceStore{
		preferences: map[string]NotificationType{
			"user001": NotificationTypeEmail,
			"user002": NotificationTypeSMS,
			"user003": NotificationTypePush,
			"user004": NotificationTypeSlack,
		},
	}
}

func (s *mockUserPreferenceStore) GetPreferredNotificationType(ctx context.Context, userID string) (NotificationType, error) {
	pref, ok := s.preferences[userID]
	if !ok {
		return "", errors.New("user preference not found")
	}
	return pref, nil
}
