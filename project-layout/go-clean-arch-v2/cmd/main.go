package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/article"
	"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/author"
	"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/pkg/config"
	"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/pkg/database"
	"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/pkg/middleware"

	"github.com/labstack/echo"

	"go.uber.org/fx"

	_ "github.com/go-sql-driver/mysql"
)

func registerHooks(lifecycle fx.Lifecycle, e *echo.Echo, cfg *config.Config) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				fmt.Println("Starting server")
				go e.Start(cfg.Server.Address)
				return nil
			},
			OnStop: func(context.Context) error {
				fmt.Println("Stopping server")
				return nil
			},
		},
	)
}

func NewEcho() *echo.Echo {
	e := echo.New()
	middle := middleware.InitMiddleware()
	e.Use(middle.CORS)
	return e
}

func ProvideBasicConfig(cfg *config.Config) time.Duration {
	return time.Duration(cfg.Context.Timeout) * time.Second
}

func main() {
	app := fx.New(
		fx.Provide(
			config.New,
			database.New,
			NewEcho,
			ProvideBasicConfig,

			article.NewArticleHandler,

			article.NewArticleUsecase,
			article.NewMysqlArticleRepository,

			author.NewMysqlAuthorRepository,
		),
		fx.Invoke(registerHooks),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	<-app.Done()
}
