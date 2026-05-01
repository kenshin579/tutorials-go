package domain

import "time"

type Action string

const (
	ActionRead Action = "read"
	ActionEdit Action = "edit"
)

func (a Action) Valid() bool {
	return a == ActionRead || a == ActionEdit
}

// ACLEntry: 한 사용자의 한 페이지에 대한 한 action 권한.
// (page_id, user_id, action) 복합 unique 인덱스로 중복 방지.
type ACLEntry struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PageID    uint      `gorm:"not null;uniqueIndex:idx_page_user_action" json:"page_id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_page_user_action" json:"user_id"`
	Action    Action    `gorm:"size:20;not null;uniqueIndex:idx_page_user_action" json:"action"`
	CreatedAt time.Time `json:"created_at"`
}
