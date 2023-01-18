package oauth

import (
	"github.com/eko/authz/backend/internal/oauth/client"
	"github.com/eko/authz/backend/internal/oauth/server"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("oauth",
		fx.Provide(
			client.NewManager,

			server.NewClientStore,
			server.NewManager,
			server.NewServer,
			server.NewTokenStore,
			server.NewAccessGenerate,

			func(store *server.ClientStore) oauth2.ClientStore { return store },
			func(manager *manage.Manager) oauth2.Manager { return manager },
			func(store *server.TokenStore) oauth2.TokenStore { return store },
		),
	)
}
