package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

// PageRepository는 GORM 기반 domain.PageRepository 구현체다.
type PageRepository struct{ db *gorm.DB }

// NewPageRepository는 *gorm.DB에서 동작하는 PageRepository를 생성한다.
func NewPageRepository(db *gorm.DB) *PageRepository { return &PageRepository{db: db} }

// Create는 새 페이지를 저장한다 (저장 후 GORM이 p.ID를 자동 채움).
func (r *PageRepository) Create(p *domain.Page) error { return r.db.Create(p).Error }

// FindByID는 id로 페이지를 조회한다. 없으면 domain.ErrNotFound를 반환한다.
func (r *PageRepository) FindByID(id uint) (*domain.Page, error) {
	var p domain.Page
	if err := r.db.First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "page"}
		}
		return nil, err
	}
	return &p, nil
}

// Update는 기존 페이지의 변경사항을 저장한다.
func (r *PageRepository) Update(p *domain.Page) error { return r.db.Save(p).Error }

// ListAccessibleBy는 본인이 owner이거나 ACLEntry로 어떤 action이든 부여받은 페이지를 반환한다.
// 같은 페이지가 owner+ACL 두 경로로 매칭되더라도 Distinct로 한 번만 노출된다.
func (r *PageRepository) ListAccessibleBy(userID uint) ([]domain.Page, error) {
	var pages []domain.Page
	err := r.db.
		Distinct("pages.*").
		Joins("LEFT JOIN acl_entries ON acl_entries.page_id = pages.id").
		Where("pages.owner_id = ? OR acl_entries.user_id = ?", userID, userID).
		Order("pages.id ASC").
		Find(&pages).Error
	return pages, err
}
