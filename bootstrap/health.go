package bootstrap

import (
	"time"

	"github.com/lengocson131002/go-clean/infras/health"
	healthchecks "github.com/lengocson131002/go-clean/pkg/health"
)

type HealthCheckerEndpoint interface {
	LivenessCheckEndpoint() healthchecks.ApplicationHealthDetailed
	ReadinessCheckEnpoint() healthchecks.ApplicationHealthDetailed
}

type healthCheckerEndpoint struct {
	healhChecker healthchecks.HealthChecker
}

func NewHealthEndpoint(cfg *ServerConfig) *healthCheckerEndpoint {
	// Init health
	healthChecker := healthchecks.NewHealthChecker(cfg.Name, cfg.AppVersion)

	// check Garbage Collector
	gcChecker := health.NewGarbageCollectionMaxChecker(time.Millisecond * time.Duration(cfg.GcPauseThresholdMs))
	healthChecker.AddLivenessCheck("garbage collector check", gcChecker)

	// check Goroutine
	grChecker := health.NewGoroutineChecker(cfg.GrRunningThreshold)
	healthChecker.AddLivenessCheck("goroutine checker", grChecker)

	// check env file
	envFileChecker := health.NewEnvChecker(cfg.EnvFilePath)
	healthChecker.AddReadinessCheck("env file checker", envFileChecker)

	// check network
	pingChecker := health.NewPingChecker("https://google.com", "GET", time.Millisecond*time.Duration(200), nil, nil)
	healthChecker.AddReadinessCheck("ping check", pingChecker)

	return &healthCheckerEndpoint{
		healhChecker: healthChecker,
	}
}

func (app healthCheckerEndpoint) LivenessCheckEndpoint() healthchecks.ApplicationHealthDetailed {
	return app.healhChecker.LivenessCheck()
}

func (app healthCheckerEndpoint) ReadinessCheckEnpoint() healthchecks.ApplicationHealthDetailed {
	return app.healhChecker.RedinessCheck()
}
