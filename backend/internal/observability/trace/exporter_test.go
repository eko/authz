package trace

import (
	"testing"

	"github.com/eko/authz/backend/configs"
	"github.com/stretchr/testify/assert"
)

func TestNewExporter(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          *configs.App
		wantExporter bool
		wantErr      error
	}{
		{
			name:         "tracing disabled",
			cfg:          &configs.App{TraceEnabled: false, TraceExporter: jaegerExporter},
			wantExporter: false,
		},
		{
			name:         "otlpgrpc exporter",
			cfg:          &configs.App{TraceEnabled: true, TraceExporter: otlpgrpcExporter, TraceOtlpEndpoint: "localhost:4317"},
			wantExporter: true,
		},
		{
			name:    "removed zipkin exporter",
			cfg:     &configs.App{TraceEnabled: true, TraceExporter: configs.TraceExporterZipkin},
			wantErr: configs.ErrTraceExporterZipkinRemoved,
		},
		{
			name:    "removed zipkin exporter with tracing disabled",
			cfg:     &configs.App{TraceEnabled: false, TraceExporter: configs.TraceExporterZipkin},
			wantErr: configs.ErrTraceExporterZipkinRemoved,
		},
		{
			name:    "unknown exporter",
			cfg:     &configs.App{TraceEnabled: true, TraceExporter: "unknown"},
			wantErr: ErrUnknownExporter,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			exporter, err := NewExporter(testCase.cfg)

			if testCase.wantErr != nil {
				assert.ErrorIs(t, err, testCase.wantErr)
				assert.Nil(t, exporter)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.wantExporter, exporter != nil)
		})
	}
}
