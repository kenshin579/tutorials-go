package clock

import "time"

// Coupon은 만료 시각을 가진 쿠폰이다.
type Coupon struct {
	Code      string
	ExpiresAt time.Time
}

// IsExpiredAt은 주어진 시각 기준으로 쿠폰 만료 여부를 반환한다.
// 시간을 파라미터로 받는 순수 함수라서 실행 시점과 무관하게 테스트할 수 있다.
func (c Coupon) IsExpiredAt(now time.Time) bool {
	return now.After(c.ExpiresAt)
}
