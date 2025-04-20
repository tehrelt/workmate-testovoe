package tracer

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	TracerKey = "tracer"
)

func newExporter(ctx context.Context, url string) (tracesdk.SpanExporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithEndpointURL(url))
}

func newProvider(exp tracesdk.SpanExporter, serviceName string) (*tracesdk.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(r),
	), nil
}

func SetupTracer(ctx context.Context, endpoint, appname string) (trace.Tracer, error) {
	slog.Debug("tracing setup",
		slog.String("endpoint", endpoint),
		slog.String("appname", appname),
	)

	exporter, err := newExporter(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	provider, err := newProvider(exporter, appname)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(provider)
	slog.Debug("tracer setup complete")

	return provider.Tracer(TracerKey), nil
}
