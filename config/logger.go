package config

import (
	"sync"

	"github.com/lengocson131002/go-clean/pkg/logger"
)

var log logger.LoggerInterface
var onceLogger sync.Once

type LoggerConfig struct {
	Level string
}

func GetLogger() logger.LoggerInterface {
	if log == nil {
		onceLogger.Do(func() {
			log = logger.NewLogrus()
		})
	}
	return log
}
