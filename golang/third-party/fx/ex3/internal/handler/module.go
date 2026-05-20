package handler

import (
	"github.com/kenshin579/tutorials-go/golang/third-party/fx/ex3/internal/handler/hello"
	"go.uber.org/fx"
)

var Module = fx.Options(
	hello.Module,
)
