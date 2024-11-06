package main

import (
	"github.com/Scusemua/go-utils/config"
	"github.com/Scusemua/go-utils/logger"
)

func main() {
	config.LogLevel = logger.LOG_LEVEL_ALL
	config.Verbose = true

	log := config.GetLogger("MyLogger ")

	log.Info("Testing, 123.")
	log.Debug("Testing, 123.")
	log.Warn("Testing, 123.")
	log.Error("Testing, 123.")
	log.Trace("Testing, 123.")
}
