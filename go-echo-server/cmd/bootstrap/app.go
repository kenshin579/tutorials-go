package bootstrap

import (
	"context"
	"fmt"

	echoRoute "github.com/kenshin579/tutorials-go/go-echo-server/echo/route/http"
	pingRoute "github.com/kenshin579/tutorials-go/go-echo-server/ping/route/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	swagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			newEcho,
		),
		fx.Invoke(
			echoRoute.NewEchoHandler,
			pingRoute.NewPingHandler,
			serve,
		),
	)
}

func newEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))
	e.Debug = true
	return e
}

func serve(
	lifecycle fx.Lifecycle,
	echo *echo.Echo,
) {
	echo.GET("/swagger/*", swagger.WrapHandler)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting echo api server.")

			go func() {
				if err := echo.Start(":80"); err != nil {
					panic(fmt.Sprintf("echo.Start() err=%v", err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping echo api server.")
			return echo.Shutdown(ctx)
		},
	})
}
