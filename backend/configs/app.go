package configs

import (
	"errors"
	"time"
)

// TraceExporterZipkin is the removed zipkin trace exporter name, kept to
// reject configurations that still use it.
const TraceExporterZipkin = "zipkin"

// ErrTraceExporterZipkinRemoved is returned when the removed zipkin trace exporter is configured.
var ErrTraceExporterZipkinRemoved = errors.New(
	"APP_TRACE_EXPORTER=zipkin is no longer supported as the OpenTelemetry Zipkin exporter is deprecated: " +
		"set APP_TRACE_EXPORTER=otlpgrpc and APP_TRACE_OTLP_ENDPOINT to an OTLP endpoint instead " +
		"(Jaeger, Zipkin with OTLP ingestion or an OpenTelemetry Collector)",
)

type App struct {
	AuditCleanDelay            time.Duration `config:"app_audit_clean_delay"`
	AuditCleanDaysToKeep       int           `config:"app_audit_clean_days_to_keep"`
	AuditFlushDelay            time.Duration `config:"app_audit_flush_delay"`
	AuditResourceKindRegex     string        `config:"app_audit_resource_kind_regex"`
	DispatcherEventChannelSize int           `config:"dispatcher_event_channel_size"`
	MetricsEnabled             bool          `config:"app_metrics_enabled"`
	StatsCleanDelay            time.Duration `config:"app_stats_clean_delay"`
	StatsCleanDaysToKeep       int           `config:"app_stats_clean_days_to_keep"`
	StatsFlushDelay            time.Duration `config:"app_stats_flush_delay"`
	StatsResourceKindRegex     string        `config:"app_stats_resource_kind_regex"`
	TraceEnabled               bool          `config:"app_trace_enabled"`
	TraceExporter              string        `config:"app_trace_exporter"`
	TraceJaegerEndpoint        string        `config:"app_trace_jaeger_endpoint"`
	TraceOtlpDialTimeout       time.Duration `config:"app_trace_otlp_dial_timeout"`
	TraceOtlpEndpoint          string        `config:"app_trace_otlp_endpoint"`
	TraceSampleRatio           float64       `config:"app_trace_sample_ratio"`
}

func newApp() *App {
	return &App{
		AuditCleanDelay:            1 * time.Hour,
		AuditCleanDaysToKeep:       7,
		AuditFlushDelay:            3 * time.Second,
		AuditResourceKindRegex:     `.*`,
		DispatcherEventChannelSize: 10000,
		MetricsEnabled:             false,
		StatsCleanDelay:            1 * time.Hour,
		StatsCleanDaysToKeep:       30,
		StatsFlushDelay:            3 * time.Second,
		StatsResourceKindRegex:     `.*`,
		TraceEnabled:               false,
		TraceExporter:              "jaeger",
		TraceJaegerEndpoint:        "localhost:14250",
		TraceOtlpDialTimeout:       3 * time.Second,
		TraceOtlpEndpoint:          "localhost:30080",
		TraceSampleRatio:           1.0,
	}
}

// validate rejects trace exporters that are no longer supported.
func (a *App) validate() error {
	if a.TraceExporter == TraceExporterZipkin {
		return ErrTraceExporterZipkinRemoved
	}

	return nil
}
