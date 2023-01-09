package grpc

import (
	"github.com/eko/authz/backend/internal/grpc/handler"
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("grpc",
		fx.Provide(
			handler.NewAuth,
			handler.NewCheck,
			handler.NewPrincipal,
			handler.NewResource,
			handler.NewPolicy,
			handler.NewRole,

			NewServer,
		),
	)
}
