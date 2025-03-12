package main

import (
	"github.com/kenshin579/tutorials-go/go-fx/ex3/internal/handler"
	"github.com/kenshin579/tutorials-go/go-fx/ex3/internal/loggerfx"
	"github.com/kenshin579/tutorials-go/go-fx/ex3/internal/routes"
	"go.uber.org/fx"
)

func main() {
	fx.New(opts()).Run()
}

func opts() fx.Option {
	return fx.Options(
		handler.Module,
		loggerfx.Module,
		fx.Invoke(routes.Register),
	)
}
