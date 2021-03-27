package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kenshin579/tutorials-go/go-fx/server"
	"go.uber.org/fx"
)

func main() {
	//1.DI 사용하지 않는 방식
	//mux := http.NewServeMux()
	//server.New(mux)
	//
	//http.ListenAndServe(":8080", mux)

	//2.DI 방식
	fx.New(
		fx.Provide(http.NewServeMux),
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
