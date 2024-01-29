package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewLogrus() *LogrusLogger {
	log := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.TraceLevel,
		Formatter: &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	}
	return &LogrusLogger{
		logger: log,
	}
}

var _ Logger = (*LogrusLogger)(nil)

// Trace implements LoggerInterface.
func (l *LogrusLogger) Trace(message string, args ...interface{}) {
	l.logger.Tracef(message, args...)
}

// Debug implements LoggerInterface.
func (l *LogrusLogger) Debug(message string, args ...interface{}) {
	l.logger.Debugf(message, args...)
}

// Error implements LoggerInterface.
func (l *LogrusLogger) Error(message string, args ...interface{}) {
	l.logger.Errorf(message, args...)
}

// Fatal implements LoggerInterface.
func (l *LogrusLogger) Fatal(message string, args ...interface{}) {
	l.logger.Fatalf(message, args...)
}

// Info implements LoggerInterface.
func (l *LogrusLogger) Info(message string, args ...interface{}) {
	l.logger.Infof(message, args...)
}

// Warn implements LoggerInterface.
func (l *LogrusLogger) Warn(message string, args ...interface{}) {
	l.logger.Warnf(message, args...)
}
