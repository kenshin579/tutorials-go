package domain

import "time"

// Confidentiality는 페이지의 분류 수준이다.
//
//	public        — 누구나 read 가능
//	internal      — 같은 부서만 read+edit 가능
//	confidential  — 같은 부서 + 정규직만 read+edit 가능
type Confidentiality string

const (
	ConfidentialityPublic       Confidentiality = "public"
	ConfidentialityInternal     Confidentiality = "internal"
	ConfidentialityConfidential Confidentiality = "confidential"
)

type Page struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	Title           string          `gorm:"size:255;not null" json:"title"`
	Content         string          `gorm:"type:text" json:"content"`
	OwnerID         uint            `gorm:"not null;index:owner_id" json:"owner_id"`
	Owner           *User           `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Confidentiality Confidentiality `gorm:"size:20;not null" json:"confidentiality"`
	// DepartmentID는 nullable: public 페이지는 부서가 없다.
	DepartmentID *uint       `gorm:"index" json:"department_id,omitempty"`
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
