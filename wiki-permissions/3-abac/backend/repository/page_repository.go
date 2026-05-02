package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
)

// PageRepository는 GORM 기반 domain.PageRepository 구현체다.
// ABAC에서는 권한 필터링을 SQL이 아닌 usecase의 정책 평가로 수행하므로 List는 모든 페이지를 반환한다.
type PageRepository struct{ db *gorm.DB }

var _ domain.PageRepository = (*PageRepository)(nil)

func NewPageRepository(db *gorm.DB) *PageRepository { return &PageRepository{db: db} }

func (r *PageRepository) Create(p *domain.Page) error { return r.db.Create(p).Error }

func (r *PageRepository) FindByID(id uint) (*domain.Page, error) {
	var p domain.Page
	if err := r.db.Preload("Department").First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "page"}
		}
		return nil, err
	}
	return &p, nil
}

func (r *PageRepository) Update(p *domain.Page) error { return r.db.Save(p).Error }

func (r *PageRepository) List() ([]domain.Page, error) {
	var pages []domain.Page
	err := r.db.Preload("Department").Order("id ASC").Find(&pages).Error
	return pages, err
}
