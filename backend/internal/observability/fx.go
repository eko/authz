package observability

import (
	"github.com/eko/authz/backend/internal/observability/metric"
	"github.com/eko/authz/backend/internal/observability/trace"
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("observability",
		fx.Provide(
			metric.NewObserver,
			metric.NewSubscriber,

			trace.NewExporter,
			trace.NewProvider,
		),

		fx.Invoke(metric.RunSubscriber),
		fx.Invoke(trace.RunProvider),
	)
}
