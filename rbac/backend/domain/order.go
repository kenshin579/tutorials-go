package domain

import (
	"errors"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

var (
	ErrInvalidTransition   = errors.New("invalid status transition")
	ErrForbiddenTransition = errors.New("forbidden status transition for this role")
)

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	ProductID  uint        `gorm:"not null" json:"product_id"`
	Product    Product     `gorm:"foreignKey:ProductID" json:"product"`
	Quantity   int         `gorm:"not null" json:"quantity"`
	TotalPrice float64     `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     OrderStatus `gorm:"size:20;not null;default:pending" json:"status"`
	OrderedBy  uint        `gorm:"not null" json:"ordered_by"`
	Orderer    User        `gorm:"foreignKey:OrderedBy" json:"orderer"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type OrderRepository interface {
	Create(order *Order) error
	FindByID(id uint) (*Order, error)
	FindAll() ([]Order, error)
	FindByUserID(userID uint) ([]Order, error)
	Update(order *Order) error
}

// validTransitions defines all valid state transitions
var validTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusPending:   {OrderStatusConfirmed, OrderStatusCancelled},
	OrderStatusConfirmed: {OrderStatusShipped, OrderStatusCancelled},
	OrderStatusShipped:   {OrderStatusCompleted},
}

// roleTransitions defines which transitions each role can perform
var roleTransitions = map[string]map[OrderStatus][]OrderStatus{
	"admin": {
		OrderStatusPending:   {OrderStatusConfirmed, OrderStatusCancelled},
		OrderStatusConfirmed: {OrderStatusShipped, OrderStatusCancelled},
		OrderStatusShipped:   {OrderStatusCompleted},
	},
	"manager": {
		OrderStatusPending:   {OrderStatusCancelled},
		OrderStatusConfirmed: {OrderStatusShipped, OrderStatusCancelled},
	},
	"user": {
		OrderStatusPending: {OrderStatusCancelled},
	},
}

// IsValidTransition checks if the transition is structurally valid
func IsValidTransition(from, to OrderStatus) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// CanTransition checks if a specific role can perform the transition
func CanTransition(role string, from, to OrderStatus) bool {
	transitions, ok := roleTransitions[role]
	if !ok {
		return false
	}
	allowed, ok := transitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// CanTransitionWithRoles checks if any of the given roles can perform the transition
func CanTransitionWithRoles(roles []string, from, to OrderStatus) bool {
	for _, role := range roles {
		if CanTransition(role, from, to) {
			return true
		}
	}
	return false
}
