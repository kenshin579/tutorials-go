package main

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func Test_basic(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.Info("hello world")
	log.Debug("hello world debug")
}
