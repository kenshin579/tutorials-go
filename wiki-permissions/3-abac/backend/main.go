package main

import (
	"log"
	"os"
	"time"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/config"
	httpx "github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/http"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/http/handler"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/repository"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/usecase"
)

func main() {
	dsn := envOr("DB_DSN", "wiki-abac.db")
	secret := envOr("JWT_SECRET", "dev-secret")
	addr := envOr("ADDR", ":8082")

	db, err := config.OpenDB(dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := config.Seed(db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	users := repository.NewUserRepository(db)
	pages := repository.NewPageRepository(db)
	depts := repository.NewDepartmentRepository(db)

	authUC := usecase.NewAuthUsecase(users, secret, 24*time.Hour)
	pageUC := usecase.NewPageUsecase(pages, users)

	deps := httpx.Deps{
		JWTSecret:  secret,
		Auth:       handler.NewAuthHandler(authUC),
		Page:       handler.NewPageHandler(pageUC),
		Department: handler.NewDepartmentHandler(depts),
	}
	e := httpx.NewRouter(deps)
	e.Logger.Fatal(e.Start(addr))
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
