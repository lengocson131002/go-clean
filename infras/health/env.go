package health

import (
	"errors"
	"fmt"
	"os"
	"time"

	healthchecks "github.com/lengocson131002/go-clean/pkg/health"
)

type EnvChecker struct {
	envFilePath string
}

// Used to check env file existed or not
func NewEnvChecker(envFilePath string) *EnvChecker {
	if envFilePath == "" {
		envFilePath = "./.env"
	}
	return &EnvChecker{
		envFilePath: envFilePath,
	}
}

// Check implements healthchecks.HealthCheckHandler.
func (ec *EnvChecker) Check(name string) healthchecks.Integration {
	var (
		start        = time.Now()
		errorMessage = ""
		status       = true
	)

	if _, err := os.Stat(ec.envFilePath); errors.Is(err, os.ErrNotExist) {
		status = false
		errorMessage = fmt.Sprintf("Env file %v does not exist", ec.envFilePath)
	}
	return healthchecks.Integration{
		Name:         name,
		Status:       status,
		ResponseTime: time.Since(start).Milliseconds(),
		Error:        errorMessage,
	}
}

var _ healthchecks.HealthCheckHandler = (*EnvChecker)(nil)
