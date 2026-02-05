package notification

import (
	"context"
	"errors"
)

// userPreferenceStore 는 UserPreferenceStore 구현체입니다.
type userPreferenceStore struct {
	preferences map[string]NotificationType
}

func NewUserPreferenceStore() UserPreferenceStore {
	return &userPreferenceStore{
		preferences: map[string]NotificationType{
			"user001": NotificationTypeEmail,
			"user002": NotificationTypeSMS,
			"user003": NotificationTypePush,
			"user004": NotificationTypeSlack,
		},
	}
}

func (s *userPreferenceStore) GetPreferredNotificationType(ctx context.Context, userID string) (NotificationType, error) {
	pref, ok := s.preferences[userID]
	if !ok {
		return "", errors.New("user preference not found")
	}
	return pref, nil
}
