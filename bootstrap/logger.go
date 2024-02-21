package bootstrap

import (
	"github.com/lengocson131002/go-clean/pkg/logger"
)

func GetLogger() logger.Logger {
	return logger.NewLogrus()
}
