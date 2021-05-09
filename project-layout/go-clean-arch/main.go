package main

import (
	"context"
	"fmt"
	"log"
	"time"

	_articleHttp "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch/article/http"
	_articleRepo "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch/article/repository/mysql"
	_articleUcase "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch/article/usecase"
	_authorRepo "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch/author/repository/mysql"
	"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch/common/config"
	"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch/common/database"

	"github.com/labstack/echo"

	"github.com/spf13/viper"

	_articleHttpMiddleware "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch/article/http/middleware"
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
				//logger.Print("Stopping admin server.")
				//return logger.Sync()
				return nil
			},
		},
	)
}

func NewEcho() *echo.Echo {
	e := echo.New()
	middle := _articleHttpMiddleware.InitMiddleware()
	e.Use(middle.CORS)
	return e
}

func ProvideBasicConfig() time.Duration {
	duration := time.Duration(viper.GetInt("context.timeout")) * time.Second
	return duration
}

//todo: 구동이후에 api 호출 주소가 등록이 안된 것 같음
func main() {
	app := fx.New(
		fx.Provide(
			config.New,
			database.New,
			NewEcho,
			ProvideBasicConfig,

			_articleHttp.NewArticleHandler,

			_articleUcase.NewArticleUsecase,
			_articleRepo.NewMysqlArticleRepository,

			_authorRepo.NewMysqlAuthorRepository,
		),
		fx.Invoke(registerHooks),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	<-app.Done()

	//아래 버전은 잘 됨
	//v := config.New()
	//db, _ := database.New(v)
	//
	//e := NewEcho()
	//
	//authorRepo := _authorRepo.NewMysqlAuthorRepository(db)
	//ar := _articleRepo.NewMysqlArticleRepository(db)
	//
	//timeoutContext := time.Duration(v.GetInt("context.timeout")) * time.Second
	//au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	//_articleHttp.NewArticleHandler(e, au)
	//
	//e.Start(v.GetString("server.address"))
}
