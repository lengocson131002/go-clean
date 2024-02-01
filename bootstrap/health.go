package bootstrap

import healthchecks "github.com/lengocson131002/go-clean/pkg/health"

type HealthCheckerEndpoint interface {
	LivenessCheckEndpoint() healthchecks.ApplicationHealthDetailed
	ReadinessCheckEnpoint() healthchecks.ApplicationHealthDetailed
}

type healthCheckerEndpoint struct {
	healhChecker healthchecks.HealthChecker
}

func NewHealthChecker() {
	// Init healt

}
