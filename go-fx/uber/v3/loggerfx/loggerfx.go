package loggerfx

import (
	"log"
	"os"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewLogger,
	),
)

type Logger interface {
	Println(v ...interface{})
}

func NewLogger() Logger {
	return log.New(os.Stdout, "[ACME]", 0)
}
