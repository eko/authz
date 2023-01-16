package audit

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("audit",
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
