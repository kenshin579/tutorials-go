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

	"github.com/spf13/viper"

	"go.uber.org/fx"

	_ "github.com/go-sql-driver/mysql"
)

func registerHooks(lifecycle fx.Lifecycle, e *echo.Echo, v *viper.Viper) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				fmt.Println("Starting server")
				go e.Start(v.GetString("server.address"))
				return nil
			},
			OnStop: func(context.Context) error {
				fmt.Println("Stopping server")
				// logger.Print("Stopping admin server.")
				// return logger.Sync()
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

func ProvideBasicConfig() time.Duration {
	duration := time.Duration(viper.GetInt("context.timeout")) * time.Second
	return duration
}

// todo: 구동이후에 api 호출 주소가 등록이 안된 것 같음
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

	// 아래 버전은 잘 됨
	// v := config.New()
	// db, _ := database.New(v)
	//
	// e := NewEcho()
	//
	// authorRepo := author.NewMysqlAuthorRepository(db)
	// ar := article.NewMysqlArticleRepository(db)
	//
	// timeoutContext := time.Duration(v.GetInt("context.timeout")) * time.Second
	// au := article.NewArticleUsecase(ar, authorRepo, timeoutContext)
	// article.NewArticleHandler(e, au)
	//
	// e.Start(v.GetString("server.address"))
}
