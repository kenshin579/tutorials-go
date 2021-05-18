package main

import (
	"log"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/cronner"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/cmd/cli"
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/config"
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/router"

	_scheduleHandler "github.com/kenshin579/tutorials-go/go-schedule/scheduler/schedule/route/http"
	_scheduleStore "github.com/kenshin579/tutorials-go/go-schedule/scheduler/schedule/store/local"
	_scheduleUsecase "github.com/kenshin579/tutorials-go/go-schedule/scheduler/schedule/usecase"
)

func main() {
	flag := cli.ParseFlags()
	cfg, err := config.New(flag.ConfigPath)
	if err != nil {
		log.Panic("config error: ", err)
	}

	c, err := cronner.New()
	if err != nil {
		log.Panic(err)
	}

	r := router.New()
	api := r.Group("/api")

	scheduleStore := _scheduleStore.NewLocalScheduleStore(c)
	scheduleUsecase := _scheduleUsecase.NewScheduleUsecase(scheduleStore, cfg)
	_scheduleHandler.NewScheduleHandler(api, scheduleUsecase)

	scheduleUsecase.InitializeJobs()
	r.Logger.Fatal(r.Start(cfg.Listen))
}
