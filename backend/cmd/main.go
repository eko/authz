package main

import (
	"context"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/compile"
	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/helper"
	"github.com/eko/authz/backend/internal/http"
	"github.com/eko/authz/backend/internal/log"
	"github.com/eko/authz/backend/internal/manager"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(context.Background),

		compile.FxModule(),
		configs.FxModule(),
		database.FxModule(),
		event.FxModule(),
		helper.FxModule(),
		http.FxModule(),
		log.FxModule(),
		manager.FxModule(),

		fx.Invoke(http.Run),
	).Run()
}
