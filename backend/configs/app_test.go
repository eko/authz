package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp_Validate(t *testing.T) {
	testCases := []struct {
		name          string
		traceEnabled  bool
		traceExporter string
		expectedErr   error
	}{
		{
			name:          "jaeger exporter is accepted",
			traceExporter: "jaeger",
		},
		{
			name:          "otlpgrpc exporter is accepted",
			traceExporter: "otlpgrpc",
		},
		{
			name:          "removed zipkin exporter is rejected when tracing is enabled",
			traceEnabled:  true,
			traceExporter: TraceExporterZipkin,
			expectedErr:   ErrTraceExporterZipkinRemoved,
		},
		{
			name:          "zipkin exporter is ignored when tracing is disabled",
			traceExporter: TraceExporterZipkin,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := newApp()
			cfg.TraceEnabled = testCase.traceEnabled
			cfg.TraceExporter = testCase.traceExporter

			err := cfg.validate()

			assert.ErrorIs(t, err, testCase.expectedErr)
		})
	}
}
