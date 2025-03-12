package main

import (
	"github.com/kenshin579/tutorials-go/third-party/go-echo-web-framework/article/handler"
	"github.com/kenshin579/tutorials-go/third-party/go-echo-web-framework/article/router"
	"github.com/kenshin579/tutorials-go/third-party/go-echo-web-framework/article/store"
	"github.com/kenshin579/tutorials-go/third-party/go-echo-web-framework/article/usecase"
)

func main() {
	router := router.New()
	v1 := router.Group("/api")

	articleStore := store.NewArticleStore()
	articleUsecase := usecase.NewArticleUsecase(articleStore)

	handler := handler.NewHandler(articleUsecase)
	handler.Register(v1)
	router.Logger.Fatal(router.Start("127.0.0.1:8080"))
}
