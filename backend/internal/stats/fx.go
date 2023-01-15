package stats

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("stats",
		fx.Provide(
			NewCleaner,
			NewSubscriber,
		),
		fx.Invoke(
			RunCleaner,
			RunSubscriber,
		),
	)
}
