package trace

import (
	"github.com/eko/authz/backend/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func NewProvider(
	cfg *configs.App,
	exporter tracesdk.SpanExporter,
) (*tracesdk.TracerProvider, error) {
	if !cfg.TraceEnabled {
		return nil, nil
	}

	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(cfg.TraceSampleRatio)),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(configs.ApplicationName),
		)),
	)

	return tracerProvider, nil
}

func RunProvider(tracerProvider *tracesdk.TracerProvider) {
	if tracerProvider == nil {
		return
	}

	otel.SetTracerProvider(tracerProvider)
}
