package security

import (
	"github.com/eko/authz/backend/internal/security/jwt"
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("security",
		fx.Provide(
			jwt.NewManager,
		),
	)
}
