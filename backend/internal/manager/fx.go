package manager

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("manager",
		fx.Provide(
			New,
			func(manager *manager) Manager { return manager },
		),
	)
}
