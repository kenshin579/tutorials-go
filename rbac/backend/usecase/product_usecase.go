package usecase

import "github.com/kenshin579/tutorials-go/rbac/backend/domain"

// ProductUsecase defines product management operations.
type ProductUsecase interface {
	List(roles []string) ([]domain.Product, error)
	GetByID(id uint) (*domain.Product, error)
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id uint) error
}

type productUsecase struct {
	productRepo domain.ProductRepository
}

// NewProductUsecase creates a new ProductUsecase.
func NewProductUsecase(productRepo domain.ProductRepository) ProductUsecase {
	return &productUsecase{productRepo: productRepo}
}

func (u *productUsecase) List(roles []string) ([]domain.Product, error) {
	// user role sees only active products
	activeOnly := !hasRole(roles, "admin") && !hasRole(roles, "manager")
	return u.productRepo.FindAll(activeOnly)
}

func (u *productUsecase) GetByID(id uint) (*domain.Product, error) {
	return u.productRepo.FindByID(id)
}

func (u *productUsecase) Create(product *domain.Product) error {
	return u.productRepo.Create(product)
}

func (u *productUsecase) Update(product *domain.Product) error {
	return u.productRepo.Update(product)
}

func (u *productUsecase) Delete(id uint) error {
	return u.productRepo.Delete(id)
}

func hasRole(roles []string, target string) bool {
	for _, r := range roles {
		if r == target {
			return true
		}
	}
	return false
}
