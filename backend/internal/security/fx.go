package security

import (
	"github.com/eko/authz/backend/internal/security/paseto"
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("security",
		fx.Provide(
			paseto.NewManager,
		),
	)
}
