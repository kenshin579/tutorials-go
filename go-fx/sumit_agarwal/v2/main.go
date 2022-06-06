package main

import (
	"context"
	"net/http"

	"github.com/kenshin579/tutorials-go/go-fx/sumit_agarwal/v2/server"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(http.NewServeMux), //NewServeMux를 server.New로 넘겨주기 위해서 추가됨
		//todo: 어떤 건 Provide로 넘기고 어떤 건 Invoke를 넘기는 이유?
		fx.Invoke(server.New),
		fx.Invoke(registerHooks),
	).Run()
}

func registerHooks(lifecycle fx.Lifecycle, mux *http.ServeMux) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go http.ListenAndServe(":8080", mux)
				return nil
			},
		},
	)
}
