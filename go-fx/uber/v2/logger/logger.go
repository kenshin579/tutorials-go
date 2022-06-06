package logger

import (
	"log"
	"os"
)

type Logger interface {
	Println(v ...interface{})
}

func NewLogger() Logger {
	return log.New(os.Stdout, "[ACME]", 0)
}
