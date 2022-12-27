package fixtures

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("fixtures",
		fx.Provide(
			NewInitializer,
		),
		fx.Invoke(
			func(initializer *initializer) error {
				return initializer.Initialize()
			},
		),
	)
}
