package trace

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type TraceConfig struct {
	ServiceName string
	Endpoint    string
	Headers     map[string]string
}

type Span interface {
	Finish()
	FinishWithOptions(opts FinishOptions)
	SetOperationName(operationName string) Span
	Tracer() Tracer
}

type FinishOptions struct {
	FinishTime time.Time
}

type Tracer interface {
	StartHttpServerTracerSpan(ctx context.Context, operationName string) Span
	StartGrpcServerTracerSpan(ctx context.Context, operationName string) Span
	StartKafkaConsumerTracerSpan(ctx context.Context, operationName string) Span
	StartDatasourceTraceSpan(ctx context.Context, operationName string) Span
}

func Init(config TraceConfig) (*sdktrace.TracerProvider, error) {
	secureOption := otlptracegrpc.WithInsecure()

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(config.Endpoint),
			otlptracegrpc.WithHeaders(config.Headers),
		),
	)

	if err != nil {
		return nil, err
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(config.ServiceName))),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return traceProvider, nil
}
