package stats

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("stats",
		fx.Provide(
			NewSubscriber,
		),
		fx.Invoke(RunSubscriber),
	)
}
