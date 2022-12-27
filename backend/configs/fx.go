package configs

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("configs",
		fx.Provide(
			Load,
			func(cfg *Base) *Database { return cfg.Database },
			func(cfg *Base) *Logger { return cfg.Logger },
			func(cfg *Base) *HTTPServer { return cfg.HTTPServer },
			func(cfg *Base) *Auth { return cfg.Auth },
			func(cfg *Base) *User { return cfg.User },
		),
	)
}
