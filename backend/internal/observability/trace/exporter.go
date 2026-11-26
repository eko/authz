package trace

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eko/authz/backend/configs"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	jaegerExporter   = "jaeger"
	otlpgrpcExporter = "otlpgrpc"
)

var (
	// ErrUnknownExporter is returned when an exporter is not implemented.
	ErrUnknownExporter = errors.New("exporter provided in configuration is unknown")
)

func NewExporter(cfg *configs.App) (tracesdk.SpanExporter, error) {
	if !cfg.TraceEnabled {
		return nil, nil
	}

	if cfg.TraceExporter == configs.TraceExporterZipkin {
		return nil, configs.ErrTraceExporterZipkinRemoved
	}

	switch cfg.TraceExporter {
	case jaegerExporter:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		conn, err := grpc.NewClient(
			cfg.TraceJaegerEndpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC connection to otlp collector: %w", err)
		}

		return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))

	case otlpgrpcExporter:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		conn, err := grpc.NewClient(
			cfg.TraceOtlpEndpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC connection to otlp collector: %w", err)
		}

		return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	}

	return nil, ErrUnknownExporter
}
