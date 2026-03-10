package config

import (
	"log"

	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AutoMigrate runs GORM auto migration for all domain types.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Permission{},
		&domain.Role{},
		&domain.User{},
		&domain.Product{},
		&domain.Order{},
	)
}

// SeedData creates initial permissions, roles, users, products and orders.
// Uses FirstOrCreate pattern to be idempotent.
func SeedData(db *gorm.DB) error {
	// 1. Permissions (17개)
	permDefs := []domain.Permission{
		{Resource: "users", Action: "create", Description: "사용자 생성"},
		{Resource: "users", Action: "read", Description: "사용자 조회"},
		{Resource: "users", Action: "update", Description: "사용자 수정"},
		{Resource: "users", Action: "delete", Description: "사용자 삭제"},
		{Resource: "roles", Action: "create", Description: "Role 생성"},
		{Resource: "roles", Action: "read", Description: "Role 조회"},
		{Resource: "roles", Action: "update", Description: "Role 수정"},
		{Resource: "roles", Action: "delete", Description: "Role 삭제"},
		{Resource: "products", Action: "create", Description: "상품 등록"},
		{Resource: "products", Action: "read", Description: "상품 조회"},
		{Resource: "products", Action: "update", Description: "상품 수정"},
		{Resource: "products", Action: "delete", Description: "상품 삭제"},
		{Resource: "products", Action: "status:update", Description: "상품 상태 변경"},
		{Resource: "orders", Action: "create", Description: "주문 생성"},
		{Resource: "orders", Action: "read", Description: "주문 조회"},
		{Resource: "orders", Action: "status:update", Description: "주문 상태 변경"},
		{Resource: "orders", Action: "cancel", Description: "주문 취소"},
	}

	permissions := make([]domain.Permission, len(permDefs))
	for i, p := range permDefs {
		result := db.Where("resource = ? AND action = ?", p.Resource, p.Action).FirstOrCreate(&permissions[i], p)
		if result.Error != nil {
			return result.Error
		}
	}

	// Build permission lookup by key
	permMap := make(map[string]domain.Permission)
	for _, p := range permissions {
		permMap[p.Key()] = p
	}

	// 2. Roles (3개)
	type roleDef struct {
		Name        string
		Description string
		PermKeys    []string
	}

	roleDefs := []roleDef{
		{
			Name:        "admin",
			Description: "관리자 - 모든 권한",
			PermKeys: []string{
				"users:create", "users:read", "users:update", "users:delete",
				"roles:create", "roles:read", "roles:update", "roles:delete",
				"products:create", "products:read", "products:update", "products:delete", "products:status:update",
				"orders:create", "orders:read", "orders:status:update", "orders:cancel",
			},
		},
		{
			Name:        "manager",
			Description: "매니저 - 제한된 관리 권한",
			PermKeys: []string{
				"users:read", "roles:read",
				"products:create", "products:read", "products:update", "products:status:update",
				"orders:create", "orders:read", "orders:status:update", "orders:cancel",
			},
		},
		{
			Name:        "user",
			Description: "일반 사용자",
			PermKeys: []string{
				"products:read",
				"orders:create", "orders:read", "orders:cancel",
			},
		},
	}

	roles := make(map[string]domain.Role)
	for _, rd := range roleDefs {
		var role domain.Role
		result := db.Where("name = ?", rd.Name).FirstOrCreate(&role, domain.Role{
			Name:        rd.Name,
			Description: rd.Description,
		})
		if result.Error != nil {
			return result.Error
		}

		// Assign permissions
		var perms []domain.Permission
		for _, key := range rd.PermKeys {
			if p, ok := permMap[key]; ok {
				perms = append(perms, p)
			}
		}
		if err := db.Model(&role).Association("Permissions").Replace(perms); err != nil {
			return err
		}
		roles[rd.Name] = role
	}

	// 3. Test users
	type userDef struct {
		Email    string
		Password string
		Name     string
		RoleName string
	}

	userDefs := []userDef{
		{Email: "admin@example.com", Password: "admin123", Name: "Admin User", RoleName: "admin"},
		{Email: "manager@example.com", Password: "manager123", Name: "Manager User", RoleName: "manager"},
		{Email: "user@example.com", Password: "user123", Name: "Normal User", RoleName: "user"},
	}

	users := make(map[string]domain.User)
	for _, ud := range userDefs {
		var user domain.User
		result := db.Where("email = ?", ud.Email).First(&user)
		if result.Error != nil {
			hash, err := bcrypt.GenerateFromPassword([]byte(ud.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			user = domain.User{
				Email:        ud.Email,
				PasswordHash: string(hash),
				Name:         ud.Name,
			}
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}

		role := roles[ud.RoleName]
		if err := db.Model(&user).Association("Roles").Replace([]domain.Role{role}); err != nil {
			return err
		}
		users[ud.RoleName] = user
	}

	// 4. Sample products
	type productDef struct {
		Name      string
		Price     float64
		Status    domain.ProductStatus
		CreatedBy uint
	}

	productDefs := []productDef{
		{Name: "노트북 Pro 15", Price: 1500000, Status: domain.ProductStatusActive, CreatedBy: users["admin"].ID},
		{Name: "무선 키보드", Price: 89000, Status: domain.ProductStatusActive, CreatedBy: users["manager"].ID},
		{Name: "USB-C 허브", Price: 45000, Status: domain.ProductStatusActive, CreatedBy: users["manager"].ID},
		{Name: "단종 마우스", Price: 25000, Status: domain.ProductStatusInactive, CreatedBy: users["admin"].ID},
	}

	products := make([]domain.Product, len(productDefs))
	for i, pd := range productDefs {
		result := db.Where("name = ?", pd.Name).FirstOrCreate(&products[i], domain.Product{
			Name:      pd.Name,
			Price:     pd.Price,
			Status:    pd.Status,
			CreatedBy: pd.CreatedBy,
		})
		if result.Error != nil {
			return result.Error
		}
	}

	// 5. Sample orders
	type orderDef struct {
		ProductIdx int
		Quantity   int
		Status     domain.OrderStatus
		OrderedBy  uint
	}

	orderDefs := []orderDef{
		{ProductIdx: 0, Quantity: 1, Status: domain.OrderStatusPending, OrderedBy: users["user"].ID},
		{ProductIdx: 1, Quantity: 2, Status: domain.OrderStatusConfirmed, OrderedBy: users["user"].ID},
		{ProductIdx: 2, Quantity: 3, Status: domain.OrderStatusShipped, OrderedBy: users["manager"].ID},
	}

	for _, od := range orderDefs {
		product := products[od.ProductIdx]
		totalPrice := product.Price * float64(od.Quantity)

		var count int64
		db.Model(&domain.Order{}).Where("product_id = ? AND ordered_by = ?", product.ID, od.OrderedBy).Count(&count)
		if count == 0 {
			order := domain.Order{
				ProductID:  product.ID,
				Quantity:   od.Quantity,
				TotalPrice: totalPrice,
				Status:     od.Status,
				OrderedBy:  od.OrderedBy,
			}
			if err := db.Create(&order).Error; err != nil {
				return err
			}
		}
	}

	log.Println("Seed data loaded successfully")
	return nil
}
