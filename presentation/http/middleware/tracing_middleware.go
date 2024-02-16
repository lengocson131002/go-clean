package middleware

import (
	"github.com/gofiber/fiber/v2"
	ot "github.com/lengocson131002/go-clean/pkg/trace/opentelemetry"
	"go.opentelemetry.io/otel/trace"
)

type TracingMiddleware struct {
	tracer *ot.OpenTelemetryTracer
}

func NewTracingMiddleware(tracer *ot.OpenTelemetryTracer) *TracingMiddleware {
	return &TracingMiddleware{tracer}
}

func (m *TracingMiddleware) Handle(ctx *fiber.Ctx) error {
	spanOpts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindServer),
	}
	_, span := m.tracer.StartSpanFromContext(ctx.UserContext(), "servicename", spanOpts...)
	defer span.End()

	return ctx.Next()
}
