package helper

import (
	"github.com/eko/authz/backend/internal/helper/time"
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("helper",
		fx.Provide(
			time.NewClock,
		),
	)
}
