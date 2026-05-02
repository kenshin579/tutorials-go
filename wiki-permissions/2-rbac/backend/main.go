package main

import (
	"log"
	"os"
	"time"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/config"
	httpx "github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/http"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/http/handler"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/repository"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/usecase"
)

func main() {
	dsn := envOr("DB_DSN", "wiki-rbac.db")
	secret := envOr("JWT_SECRET", "dev-secret")
	addr := envOr("ADDR", ":8081")

	db, err := config.OpenDB(dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := config.Seed(db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	users := repository.NewUserRepository(db)
	pages := repository.NewPageRepository(db)
	roles := repository.NewRoleRepository(db)
	perms := repository.NewPermissionRepository(db)

	authUC := usecase.NewAuthUsecase(users, perms, secret, 24*time.Hour)
	pageUC := usecase.NewPageUsecase(pages, perms)
	roleUC := usecase.NewRoleUsecase(users, roles, perms)

	deps := httpx.Deps{
		JWTSecret: secret,
		Auth:      handler.NewAuthHandler(authUC),
		Page:      handler.NewPageHandler(pageUC),
		Role:      handler.NewRoleHandler(roleUC),
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
