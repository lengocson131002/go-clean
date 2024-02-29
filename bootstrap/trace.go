package bootstrap

import (
	"context"
	"fmt"

	"github.com/lengocson131002/go-clean/pkg/config"
	"github.com/lengocson131002/go-clean/pkg/trace"
	"github.com/lengocson131002/go-clean/pkg/trace/otel"
)

func ConfigTracing(srvCfg *ServerConfig, cfg config.Configure) {
	endpoint := cfg.GetString("TRACE_ENDPOINT")
	otel.SetGlobalTracer(context.Background(), srvCfg.Name, endpoint)
}

func GetTracer(srvCfg *ServerConfig, cfg config.Configure) trace.Tracer {
	endpoint := cfg.GetString("TRACE_ENDPOINT")
	tracer, err := otel.NewOpenTelemetryTracer(context.Background(), srvCfg.Name, endpoint)
	if err != nil {
		panic(fmt.Errorf("Failed to create tracer object: %w", err))
	}
	return tracer
}
