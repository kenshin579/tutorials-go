package domain

import "time"

// EmploymentType은 사용자의 고용 형태다.
// confidential 페이지 접근 정책에서 정규직(fulltime)만 허용하는 조건으로 사용된다.
type EmploymentType string

const (
	EmploymentFulltime EmploymentType = "fulltime"
	EmploymentContract EmploymentType = "contract"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Email          string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	PasswordHash   string         `gorm:"size:255;not null" json:"-"`
	DepartmentID   uint           `gorm:"not null;index" json:"department_id"`
	Department     *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	EmploymentType EmploymentType `gorm:"size:20;not null" json:"employment_type"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
