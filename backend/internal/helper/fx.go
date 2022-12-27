package helper

import (
	"github.com/eko/authz/backend/internal/helper/time"
	"github.com/eko/authz/backend/internal/helper/token"
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("helper",
		fx.Provide(
			time.NewClock,
			token.NewGenerator,
		),
	)
}
