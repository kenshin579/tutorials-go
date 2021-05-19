package main

import (
	"log"

	_blackBoardStore "github.com/kenshin579/tutorials-go/go-redis/blackboard/blackboard/store/redis"

	_blackBoardUsecase "github.com/kenshin579/tutorials-go/go-redis/blackboard/blackboard/usecase"

	_blackBoardHandler "github.com/kenshin579/tutorials-go/go-redis/blackboard/blackboard/route/http"

	"github.com/kenshin579/tutorials-go/go-redis/blackboard/cmd/cli"
	"github.com/kenshin579/tutorials-go/go-redis/blackboard/common/config"
	"github.com/kenshin579/tutorials-go/go-redis/blackboard/common/router"
)

func main() {
	flag := cli.ParseFlags()
	cfg, err := config.New(flag.ConfigPath)

	if err != nil {
		log.Panic("config error: ", err)
	}

	r := router.New()
	g := r.Group("/api")

	blackBoardStore := _blackBoardStore.NewRedisBlackBoardStore(cfg)
	blackBoardUsecase := _blackBoardUsecase.NewBlackBoardUsecase(blackBoardStore)
	_blackBoardHandler.NewBlackBoardHandler(g, blackBoardUsecase)
	r.Logger.Fatal(r.Start(cfg.Listen))
}
