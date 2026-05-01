package main

import (
	"log"
	"os"
	"time"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/config"
	httpx "github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/http"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/http/handler"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/repository"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/usecase"
)

func main() {
	dsn := envOr("DB_DSN", "wiki-acl.db")
	secret := envOr("JWT_SECRET", "dev-secret")
	addr := envOr("ADDR", ":8080")

	db, err := config.OpenDB(dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := config.Seed(db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	users := repository.NewUserRepository(db)
	pages := repository.NewPageRepository(db)
	acls := repository.NewACLRepository(db)

	authUC := usecase.NewAuthUsecase(users, secret, 24*time.Hour)
	pageUC := usecase.NewPageUsecase(pages, acls)
	aclUC := usecase.NewACLUsecase(pages, acls)

	deps := httpx.Deps{
		JWTSecret: secret,
		Auth:      handler.NewAuthHandler(authUC),
		Page:      handler.NewPageHandler(pageUC),
		ACL:       handler.NewACLHandler(aclUC),
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
