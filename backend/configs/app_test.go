package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp_Validate(t *testing.T) {
	testCases := []struct {
		name          string
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
			name:          "removed zipkin exporter is rejected",
			traceExporter: TraceExporterZipkin,
			expectedErr:   ErrTraceExporterZipkinRemoved,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := newApp()
			cfg.TraceExporter = testCase.traceExporter

			err := cfg.validate()

			assert.ErrorIs(t, err, testCase.expectedErr)
		})
	}
}
