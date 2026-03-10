package domain

import "time"

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
)

type Product struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Name      string        `gorm:"not null" json:"name"`
	Price     float64       `gorm:"type:decimal(10,2);not null" json:"price"`
	Status    ProductStatus `gorm:"not null;default:active" json:"status"`
	CreatedBy uint          `gorm:"not null" json:"created_by"`
	Creator   User          `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type ProductRepository interface {
	Create(product *Product) error
	FindByID(id uint) (*Product, error)
	FindAll(activeOnly bool) ([]Product, error)
	Update(product *Product) error
	Delete(id uint) error
}
