package oauth

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("oauth",
		fx.Provide(
			NewClientStore,
			NewManager,
			NewServer,
			NewTokenStore,
			NewAccessGenerate,

			func(store *ClientStore) oauth2.ClientStore { return store },
			func(manager *manage.Manager) oauth2.Manager { return manager },
			func(store *TokenStore) oauth2.TokenStore { return store },
		),
	)
}
