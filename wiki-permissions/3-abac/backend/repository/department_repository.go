package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
)

// DepartmentRepository는 GORM 기반 domain.DepartmentRepository 구현체다.
type DepartmentRepository struct{ db *gorm.DB }

var _ domain.DepartmentRepository = (*DepartmentRepository)(nil)

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) FindByID(id uint) (*domain.Department, error) {
	var d domain.Department
	if err := r.db.First(&d, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "department"}
		}
		return nil, err
	}
	return &d, nil
}

func (r *DepartmentRepository) List() ([]domain.Department, error) {
	var ds []domain.Department
	err := r.db.Order("id ASC").Find(&ds).Error
	return ds, err
}
