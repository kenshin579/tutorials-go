package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kenshin579/tutorials-go/rbac/backend/config"
	apphttp "github.com/kenshin579/tutorials-go/rbac/backend/http"
	"github.com/kenshin579/tutorials-go/rbac/backend/http/handler"
	"github.com/kenshin579/tutorials-go/rbac/backend/repository"
	"github.com/kenshin579/tutorials-go/rbac/backend/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect to MySQL
	db, err := gorm.Open(mysql.Open(cfg.DB.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto migrate and seed
	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	if err := config.SeedData(db); err != nil {
		log.Fatalf("failed to seed data: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, roleRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	rbacUsecase := usecase.NewRbacUsecase(roleRepo, permRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, productRepo)

	// Handlers
	handlers := &apphttp.Handlers{
		Auth:    handler.NewAuthHandler(authUsecase, permRepo, cfg.JWT.Secret),
		User:    handler.NewUserHandler(userUsecase),
		Rbac:    handler.NewRbacHandler(rbacUsecase),
		Product: handler.NewProductHandler(productUsecase),
		Order:   handler.NewOrderHandler(orderUsecase),
	}

	// Setup Echo
	e := echo.New()

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Logger and Recover
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Setup routes
	apphttp.SetupRoutes(e, handlers, cfg.JWT.Secret, permRepo, db)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("Server starting on %s", addr)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
}
