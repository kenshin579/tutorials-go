package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

// ACLRepository는 GORM 기반 domain.ACLRepository의 부분 구현체다.
// 현재 Grant만 지원하며, FindByPageAndUser/ListByPage/Revoke는 Task 10에서 완성된다.
type ACLRepository struct{ db *gorm.DB }

// NewACLRepository는 *gorm.DB에서 동작하는 ACLRepository를 생성한다.
func NewACLRepository(db *gorm.DB) *ACLRepository { return &ACLRepository{db: db} }

// Grant는 (page, user, action) 권한을 추가한다. 동일 행이 이미 있으면 무시한다.
func (r *ACLRepository) Grant(pageID, userID uint, action domain.Action) error {
	entry := domain.ACLEntry{PageID: pageID, UserID: userID, Action: action}
	return r.db.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&entry).Error
}
