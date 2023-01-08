package database

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("database",
		fx.Provide(
			New,
			NewTransactionManager,
		),
	)
}
