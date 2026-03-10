package usecase

import (
	"errors"

	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
)

var (
	ErrForbiddenTransition = errors.New("forbidden status transition for this role")
	ErrProductNotFound     = errors.New("product not found")
)

// OrderUsecase defines order management operations.
type OrderUsecase interface {
	Create(order *domain.Order) error
	List(userID uint, roles []string) ([]domain.Order, error)
	GetByID(id uint) (*domain.Order, error)
	UpdateStatus(orderID uint, newStatus domain.OrderStatus, roles []string) error
	Cancel(orderID uint, roles []string) error
}

type orderUsecase struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

// NewOrderUsecase creates a new OrderUsecase.
func NewOrderUsecase(orderRepo domain.OrderRepository, productRepo domain.ProductRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (u *orderUsecase) Create(order *domain.Order) error {
	// Look up product to calculate total price
	product, err := u.productRepo.FindByID(order.ProductID)
	if err != nil {
		return ErrProductNotFound
	}

	order.TotalPrice = product.Price * float64(order.Quantity)
	order.Status = domain.OrderStatusPending
	return u.orderRepo.Create(order)
}

func (u *orderUsecase) List(userID uint, roles []string) ([]domain.Order, error) {
	// admin/manager see all orders, user sees own only
	if hasRole(roles, "admin") || hasRole(roles, "manager") {
		return u.orderRepo.FindAll()
	}
	return u.orderRepo.FindByUserID(userID)
}

func (u *orderUsecase) GetByID(id uint) (*domain.Order, error) {
	return u.orderRepo.FindByID(id)
}

func (u *orderUsecase) UpdateStatus(orderID uint, newStatus domain.OrderStatus, roles []string) error {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	if !domain.CanTransitionWithRoles(roles, order.Status, newStatus) {
		return ErrForbiddenTransition
	}

	order.Status = newStatus
	return u.orderRepo.Update(order)
}

func (u *orderUsecase) Cancel(orderID uint, roles []string) error {
	return u.UpdateStatus(orderID, domain.OrderStatusCancelled, roles)
}
