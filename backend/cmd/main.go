package main

import (
	"context"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/compile"
	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/entity"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/fixtures"
	internal_fx "github.com/eko/authz/backend/internal/fx"
	"github.com/eko/authz/backend/internal/grpc"
	"github.com/eko/authz/backend/internal/helper"
	"github.com/eko/authz/backend/internal/http"
	"github.com/eko/authz/backend/internal/log"
	"github.com/eko/authz/backend/internal/oauth"
	"github.com/eko/authz/backend/internal/security"
	"github.com/eko/authz/backend/internal/stats"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(context.Background),
		internal_fx.Logger,

		compile.FxModule(),
		configs.FxModule(),
		database.FxModule(),
		event.FxModule(),
		fixtures.FxModule(),
		grpc.FxModule(),
		helper.FxModule(),
		http.FxModule(),
		log.FxModule(),
		entity.FxModule(),
		oauth.FxModule(),
		security.FxModule(),
		stats.FxModule(),

		fx.Invoke(
			grpc.Run,
			http.Run,
		),
	).Run()
}
