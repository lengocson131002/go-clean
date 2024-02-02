package health

import (
	"fmt"
	"runtime"
	"time"

	healthchecks "github.com/lengocson131002/go-clean/pkg/health"
)

const (
	DEFAULT_GR_THRESHOLD = 100
)

type GoroutineChecker struct {
	threshold int
}

func NewGoroutineChecker(threshold int) *GoroutineChecker {
	if threshold == 0 {
		threshold = DEFAULT_GR_THRESHOLD
	}
	return &GoroutineChecker{
		threshold: threshold,
	}
}

// GoroutineCountCheck returns a Check that fails if too many goroutines are
// running (which could indicate a resource leak).
func (gr *GoroutineChecker) Check(name string) healthchecks.Integration {
	var (
		start        = time.Now()
		grStatus     = true
		errorMessage = ""
	)

	count := runtime.NumGoroutine()
	if count > gr.threshold {
		grStatus = false
		errorMessage = fmt.Sprintf("too many goroutines (%d > %d)", count, gr.threshold)
	}

	return healthchecks.Integration{
		Name:         name,
		Status:       grStatus,
		ResponseTime: time.Since(start).Microseconds(),
		Error:        errorMessage,
	}
}

var _ healthchecks.HealthCheckHandler = (*GoroutineChecker)(nil)
