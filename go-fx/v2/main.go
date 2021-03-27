package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kenshin579/tutorials-go/go-fx/server"
	"go.uber.org/fx"
)

func main() {
	//2.DI 방식
	fx.New(
		fx.Provide(http.NewServeMux), //NewServeMux를 server.New에 넣기 위해서 추가함
		fx.Invoke(server.New),
		fx.Invoke(registerHooks),
	).Run()

}

func registerHooks(lifecycle fx.Lifecycle, mux *http.ServeMux) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Println("starting server...")
				go http.ListenAndServe(":8080", mux)
				return nil
			},
		},
	)
}
