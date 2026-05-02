package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

// PageRepository는 GORM 기반 domain.PageRepository 구현체다.
// 1편(ACL)과 달리 ListAccessibleBy 대신 List만 제공한다 — RBAC에서는 모든 페이지를 노출하고
// 접근 가능 여부는 usecase 단에서 RBAC 평가로 결정한다.
type PageRepository struct{ db *gorm.DB }

var _ domain.PageRepository = (*PageRepository)(nil)

// NewPageRepository는 *gorm.DB에서 동작하는 PageRepository를 생성한다.
func NewPageRepository(db *gorm.DB) *PageRepository { return &PageRepository{db: db} }

// Create는 새 페이지를 저장한다.
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

// Delete는 id로 페이지를 삭제한다.
func (r *PageRepository) Delete(id uint) error { return r.db.Delete(&domain.Page{}, id).Error }

// List는 모든 페이지를 id 오름차순으로 반환한다.
func (r *PageRepository) List() ([]domain.Page, error) {
	var pages []domain.Page
	err := r.db.Order("id ASC").Find(&pages).Error
	return pages, err
}
