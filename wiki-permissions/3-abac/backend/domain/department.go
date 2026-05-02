package domain

import "time"

// Department는 조직의 부서를 표현한다 (예: Engineering, Marketing).
// ABAC에서 user.DepartmentID와 page.DepartmentID를 매칭해 internal/confidential 정책을 평가한다.
type Department struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
