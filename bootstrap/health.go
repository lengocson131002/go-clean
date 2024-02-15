package bootstrap

import (
	"time"

	health "github.com/lengocson131002/go-clean/pkg/health"
)

type HealthCheckerEndpoint interface {
	LivenessCheckEndpoint() health.ApplicationHealthDetailed
	ReadinessCheckEnpoint() health.ApplicationHealthDetailed
}

type healthCheckerEndpoint struct {
	healhChecker health.HealthChecker
}

func NewHealthEndpoint(cfg *ServerConfig) *healthCheckerEndpoint {
	// Init health
	healthChecker := health.NewHealthChecker(cfg.Name, cfg.AppVersion)

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

func (app healthCheckerEndpoint) LivenessCheckEndpoint() health.ApplicationHealthDetailed {
	return app.healhChecker.LivenessCheck()
}

func (app healthCheckerEndpoint) ReadinessCheckEnpoint() health.ApplicationHealthDetailed {
	return app.healhChecker.RedinessCheck()
}
