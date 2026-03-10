package repository

import (
	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Product").Preload("Orderer").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindAll() ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Product").Preload("Orderer").Order("id DESC").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Product").Preload("Orderer").
		Where("ordered_by = ?", userID).
		Order("id DESC").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}
