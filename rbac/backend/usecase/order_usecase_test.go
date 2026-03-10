package usecase

import (
	"errors"
	"testing"

	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock OrderRepository ---

type mockOrderRepo struct {
	orders    map[uint]*domain.Order
	nextID    uint
	createErr error
	findErr   error
	updateErr error
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{
		orders: make(map[uint]*domain.Order),
		nextID: 1,
	}
}

func (m *mockOrderRepo) Create(order *domain.Order) error {
	if m.createErr != nil {
		return m.createErr
	}
	order.ID = m.nextID
	m.nextID++
	cp := *order
	m.orders[cp.ID] = &cp
	return nil
}

func (m *mockOrderRepo) FindByID(id uint) (*domain.Order, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	order, ok := m.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	cp := *order
	return &cp, nil
}

func (m *mockOrderRepo) FindAll() ([]domain.Order, error) {
	result := make([]domain.Order, 0, len(m.orders))
	for _, o := range m.orders {
		result = append(result, *o)
	}
	return result, nil
}

func (m *mockOrderRepo) FindByUserID(userID uint) ([]domain.Order, error) {
	var result []domain.Order
	for _, o := range m.orders {
		if o.OrderedBy == userID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (m *mockOrderRepo) Update(order *domain.Order) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	cp := *order
	m.orders[order.ID] = &cp
	return nil
}

// --- Mock ProductRepository ---

type mockProductRepo struct {
	products  map[uint]*domain.Product
	findErr   error
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{
		products: make(map[uint]*domain.Product),
	}
}

func (m *mockProductRepo) Create(product *domain.Product) error {
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepo) FindByID(id uint) (*domain.Product, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	product, ok := m.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	cp := *product
	return &cp, nil
}

func (m *mockProductRepo) FindAll(activeOnly bool) ([]domain.Product, error) {
	var result []domain.Product
	for _, p := range m.products {
		if !activeOnly || p.Status == domain.ProductStatusActive {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockProductRepo) Update(product *domain.Product) error {
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepo) Delete(id uint) error {
	delete(m.products, id)
	return nil
}

// --- Helpers ---

func seedProduct(repo *mockProductRepo, id uint, price float64) *domain.Product {
	p := &domain.Product{
		ID:    id,
		Name:  "Test Product",
		Price: price,
	}
	repo.products[id] = p
	return p
}

func seedOrder(repo *mockOrderRepo, id uint, userID uint, status domain.OrderStatus) *domain.Order {
	o := &domain.Order{
		ID:        id,
		ProductID: 1,
		Quantity:  1,
		OrderedBy: userID,
		Status:    status,
	}
	repo.orders[id] = o
	return o
}

// --- Tests ---

func TestOrderUsecase_Create_CalculatesTotalPrice(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	seedProduct(productRepo, 1, 50.0)

	uc := NewOrderUsecase(orderRepo, productRepo)

	order := &domain.Order{
		ProductID: 1,
		Quantity:  3,
		OrderedBy: 10,
	}

	err := uc.Create(order)
	require.NoError(t, err)
	assert.Equal(t, 150.0, order.TotalPrice)
	assert.Equal(t, domain.OrderStatusPending, order.Status)
}

func TestOrderUsecase_Create_ProductNotFound_ReturnsError(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	// no products seeded

	uc := NewOrderUsecase(orderRepo, productRepo)

	order := &domain.Order{
		ProductID: 99,
		Quantity:  1,
		OrderedBy: 10,
	}

	err := uc.Create(order)
	assert.ErrorIs(t, err, ErrProductNotFound)
}

func TestOrderUsecase_List_AdminSeesAllOrders(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()

	seedOrder(orderRepo, 1, 10, domain.OrderStatusPending)
	seedOrder(orderRepo, 2, 20, domain.OrderStatusPending)
	seedOrder(orderRepo, 3, 30, domain.OrderStatusConfirmed)

	uc := NewOrderUsecase(orderRepo, productRepo)

	orders, err := uc.List(10, []string{"admin"})
	require.NoError(t, err)
	assert.Len(t, orders, 3)
}

func TestOrderUsecase_List_UserSeesOnlyOwnOrders(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()

	const userID uint = 10
	seedOrder(orderRepo, 1, userID, domain.OrderStatusPending)
	seedOrder(orderRepo, 2, userID, domain.OrderStatusConfirmed)
	seedOrder(orderRepo, 3, 99, domain.OrderStatusPending) // belongs to another user

	uc := NewOrderUsecase(orderRepo, productRepo)

	orders, err := uc.List(userID, []string{"user"})
	require.NoError(t, err)
	assert.Len(t, orders, 2)
	for _, o := range orders {
		assert.Equal(t, userID, o.OrderedBy)
	}
}

func TestOrderUsecase_UpdateStatus_AdminCanConfirmPendingOrder(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()

	seedOrder(orderRepo, 1, 10, domain.OrderStatusPending)

	uc := NewOrderUsecase(orderRepo, productRepo)

	err := uc.UpdateStatus(1, domain.OrderStatusConfirmed, []string{"admin"})
	require.NoError(t, err)

	updated, err := orderRepo.FindByID(1)
	require.NoError(t, err)
	assert.Equal(t, domain.OrderStatusConfirmed, updated.Status)
}

func TestOrderUsecase_UpdateStatus_ManagerCanShipConfirmedOrder(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()

	seedOrder(orderRepo, 1, 10, domain.OrderStatusConfirmed)

	uc := NewOrderUsecase(orderRepo, productRepo)

	err := uc.UpdateStatus(1, domain.OrderStatusShipped, []string{"manager"})
	require.NoError(t, err)

	updated, err := orderRepo.FindByID(1)
	require.NoError(t, err)
	assert.Equal(t, domain.OrderStatusShipped, updated.Status)
}

func TestOrderUsecase_UpdateStatus_UserCannotConfirmOrder(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()

	seedOrder(orderRepo, 1, 10, domain.OrderStatusPending)

	uc := NewOrderUsecase(orderRepo, productRepo)

	err := uc.UpdateStatus(1, domain.OrderStatusConfirmed, []string{"user"})
	assert.ErrorIs(t, err, ErrForbiddenTransition)
}

func TestOrderUsecase_Cancel_UserCanCancelPendingOrder(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()

	seedOrder(orderRepo, 1, 10, domain.OrderStatusPending)

	uc := NewOrderUsecase(orderRepo, productRepo)

	err := uc.Cancel(1, []string{"user"})
	require.NoError(t, err)

	updated, err := orderRepo.FindByID(1)
	require.NoError(t, err)
	assert.Equal(t, domain.OrderStatusCancelled, updated.Status)
}

func TestOrderUsecase_Cancel_UserCannotCancelConfirmedOrder(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()

	seedOrder(orderRepo, 1, 10, domain.OrderStatusConfirmed)

	uc := NewOrderUsecase(orderRepo, productRepo)

	err := uc.Cancel(1, []string{"user"})
	assert.ErrorIs(t, err, ErrForbiddenTransition)
}
