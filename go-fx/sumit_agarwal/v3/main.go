package main

import (
	"context"
	"net/http"

	"github.com/kenshin579/tutorials-go/go-fx/sumit_agarwal/v3/loggerfx"
	"github.com/kenshin579/tutorials-go/go-fx/sumit_agarwal/v3/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(http.NewServeMux),
		fx.Invoke(server.New),
		fx.Invoke(registerHooks),
		loggerfx.Module,
	).Run()
}

func registerHooks(
	lifecycle fx.Lifecycle, mux *http.ServeMux, logger *zap.SugaredLogger,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info("OnStart :: Listening on localhost:8080")
				go http.ListenAndServe(":8080", mux)
				return nil
			},
			OnStop: func(context.Context) error {
				logger.Info("OnStop...")
				return logger.Sync()
			},
		},
	)
}
