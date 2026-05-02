package domain

import "time"

// Permission은 한 리소스에 대한 한 액션을 표현한다 (예: pages:read).
// (resource, action) 복합 unique 인덱스로 중복 방지.
type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Resource  string    `gorm:"size:100;not null;uniqueIndex:idx_resource_action" json:"resource"`
	Action    string    `gorm:"size:100;not null;uniqueIndex:idx_resource_action" json:"action"`
	CreatedAt time.Time `json:"created_at"`
}

// Key는 "resource:action" 형태 문자열이다 (예: "pages:edit").
// 권한 평가 시 lookup 키로 사용한다.
func (p Permission) Key() string { return p.Resource + ":" + p.Action }
