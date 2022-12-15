package compile

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("compile",
		fx.Provide(
			NewCompiler,
			NewSubscriber,
			func(compiler *compiler) Compiler { return compiler },
		),
		fx.Invoke(RunSubscriber),
	)
}
