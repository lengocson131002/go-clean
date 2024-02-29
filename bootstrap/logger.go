package bootstrap

import (
	"github.com/lengocson131002/go-clean/pkg/config"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/logger/logrus"
	"github.com/lengocson131002/go-clean/pkg/trace"
)

func GetLogger(c config.Configure, tracer trace.Tracer) logger.Logger {
	levelStr := c.GetString("LOG_LEVEL")
	level, err := logger.GetLevel(levelStr)
	if err != nil {
		level = logger.InfoLevel
	}
	return logrus.NewLogrusLogger(
		logger.WithLevel(level),
		logger.WithTracer(tracer),
	)
}
