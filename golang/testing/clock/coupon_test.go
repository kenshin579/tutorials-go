package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Coupon_IsExpiredAt_경계값(t *testing.T) {
	expiresAt := time.Date(2026, 8, 12, 0, 0, 0, 0, time.UTC)
	coupon := Coupon{Code: "WELCOME10", ExpiresAt: expiresAt}

	tests := []struct {
		name    string
		now     time.Time
		expired bool
	}{
		{"만료 1초 전", expiresAt.Add(-time.Second), false},
		{"만료 시각과 동일", expiresAt, false},
		{"만료 1초 후", expiresAt.Add(time.Second), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expired, coupon.IsExpiredAt(tt.now))
		})
	}
}
