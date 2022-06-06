package main

import (
	"github.com/kenshin579/tutorials-go/go-fx/uber/v3/handlers"
	"github.com/kenshin579/tutorials-go/go-fx/uber/v3/loggerfx"
	"github.com/kenshin579/tutorials-go/go-fx/uber/v3/serverfx"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		serverfx.Module,
		loggerfx.Module,
		fx.Provide(
			handlers.NewHandler,
		),
	)
	app.Run()
}
