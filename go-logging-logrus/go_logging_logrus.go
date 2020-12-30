package main

import log "github.com/sirupsen/logrus"

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("hello world")
	log.Debug("hello world debug")
}
