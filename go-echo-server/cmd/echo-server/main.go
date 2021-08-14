package main

import (
	"github.com/kenshin579/tutorials-go/go-echo-server/common/router"
	_echoHandler "github.com/kenshin579/tutorials-go/go-echo-server/echo/route/http"
)

func main() {
	r := router.New()
	_echoHandler.NewEchoHandler(r)

	r.Logger.Fatal(r.Start(":80"))
}
