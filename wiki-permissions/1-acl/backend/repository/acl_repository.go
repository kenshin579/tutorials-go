package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

// ACLRepository는 GORM 기반 domain.ACLRepository 구현체다.
type ACLRepository struct{ db *gorm.DB }

var _ domain.ACLRepository = (*ACLRepository)(nil)

// NewACLRepository는 *gorm.DB에서 동작하는 ACLRepository를 생성한다.
func NewACLRepository(db *gorm.DB) *ACLRepository { return &ACLRepository{db: db} }

// FindByPageAndUser는 한 페이지에 대해 특정 사용자가 가진 모든 ACL 항목을 반환한다.
// (예: read와 edit을 둘 다 받은 경우 두 entry가 반환된다.)
func (r *ACLRepository) FindByPageAndUser(pageID, userID uint) ([]domain.ACLEntry, error) {
	var entries []domain.ACLEntry
	err := r.db.Where("page_id = ? AND user_id = ?", pageID, userID).Find(&entries).Error
	return entries, err
}

// ListByPage는 한 페이지의 모든 ACL 항목을 user_id, action 순으로 정렬해 반환한다.
// 페이지 owner가 공유 목록 UI에서 사용한다.
func (r *ACLRepository) ListByPage(pageID uint) ([]domain.ACLEntry, error) {
	var entries []domain.ACLEntry
	err := r.db.Where("page_id = ?", pageID).Order("user_id, action").Find(&entries).Error
	return entries, err
}

// Grant는 (page, user, action) 권한을 추가한다. 동일 행이 이미 있으면 무시한다.
func (r *ACLRepository) Grant(pageID, userID uint, action domain.Action) error {
	entry := domain.ACLEntry{PageID: pageID, UserID: userID, Action: action}
	return r.db.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&entry).Error
}

// Revoke는 (page, user)에 대해 부여된 모든 action ACL을 삭제한다.
// 한 사용자의 read/edit 권한을 함께 회수하는 시맨틱이며, 부분 회수가 필요하면 Grant로 다시 부여한다.
func (r *ACLRepository) Revoke(pageID, userID uint) error {
	return r.db.Where("page_id = ? AND user_id = ?", pageID, userID).Delete(&domain.ACLEntry{}).Error
}
