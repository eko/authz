package event

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("event",
		fx.Provide(
			NewDispatcher,
			func(dispatcher *dispatcher) Dispatcher { return dispatcher },
		),
	)
}
