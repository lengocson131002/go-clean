package bootstrap

import (
	ot "github.com/lengocson131002/go-clean/pkg/trace/opentelemetry"
)

type TraceConfig struct {
	ServiceName string `json:"serviceName"`
	Endpoint    string `json:"endpoint"`
}

func GetTracer(cfg *TraceConfig) (*ot.OpenTelemetryTracer, error) {
	return ot.NewTracer(
		ot.WithEndpoint(cfg.Endpoint),
		ot.WithServiceName(cfg.ServiceName),
	)
}
