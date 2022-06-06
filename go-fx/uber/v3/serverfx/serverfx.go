package serverfx

import (
	"context"
	"net/http"

	"go.uber.org/fx"
)

var Module = fx.Options(
	/*
		Invoke는 fx의 lifeCycle이 시작되기전에 실행된다
	*/

	fx.Invoke(
		RegisterHandlers,
		InitServer,
	),
)

func RegisterHandlers(handler http.Handler) {
	http.Handle("/", handler)
}

func InitServer(lifecycle fx.Lifecycle) {
	server := &http.Server{
		Addr: ":8080",
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Close()
		},
	})
}
