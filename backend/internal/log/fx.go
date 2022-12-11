package log

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("log",
		fx.Provide(
			New,
		),
	)
}
