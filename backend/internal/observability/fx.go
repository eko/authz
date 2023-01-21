package observability

import (
	"github.com/eko/authz/backend/internal/observability/metric"
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("observability",
		fx.Provide(
			metric.NewObserver,
			metric.NewSubscriber,
		),

		fx.Invoke(metric.RunSubscriber),
	)
}
