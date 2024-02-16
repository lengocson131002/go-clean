package opentelemetry

import (
	"context"
	"strings"

	"github.com/lengocson131002/go-clean/pkg/metadata"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	instrumentationName = "github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
)

var (
	caser = cases.Title(language.English, cases.NoLower)
)

type OpenTelemetryTracer struct {
	traceProvider *tracesdk.TracerProvider
}

func NewTracer(options ...Option) (*OpenTelemetryTracer, error) {
	var opts Options
	for _, opt := range options {
		opt(&opts)
	}

	traceProvider, err := NewTraceProvider(opts.ServiceName, opts.Endpoint)
	if err != nil {
		return nil, err
	}

	return &OpenTelemetryTracer{traceProvider: traceProvider}, nil
}

func (ot *OpenTelemetryTracer) StartSpanFromContext(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(metadata.Metadata)
	}
	propagator, carrier := otel.GetTextMapPropagator(), make(propagation.MapCarrier)
	for k, v := range md {
		for _, f := range propagator.Fields() {
			if strings.EqualFold(k, f) {
				carrier[f] = v
			}
		}
	}
	ctx = propagator.Extract(ctx, carrier)
	spanCtx := trace.SpanContextFromContext(ctx)
	ctx = baggage.ContextWithBaggage(ctx, baggage.FromContext(ctx))

	var tracer trace.Tracer
	var span trace.Span
	tp := ot.traceProvider
	if tp != nil {
		tracer = tp.Tracer(instrumentationName)
	} else {
		tracer = otel.Tracer(instrumentationName)
	}
	ctx, span = tracer.Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), name, opts...)

	carrier = make(propagation.MapCarrier)
	propagator.Inject(ctx, carrier)
	for k, v := range carrier {
		//lint:ignore SA1019 no unicode punctution handle needed
		md.Set(caser.String(k), v)
	}
	ctx = metadata.NewContext(ctx, md)

	return ctx, span
}

func (ot *OpenTelemetryTracer) Shutdown(ctx context.Context) error {
	return ot.traceProvider.Shutdown(ctx)
}
