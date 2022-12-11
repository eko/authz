package main

import (
	"context"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/http"
	"github.com/eko/authz/backend/internal/log"
	"github.com/eko/authz/backend/internal/manager"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(context.Background),

		configs.FxModule(),
		database.FxModule(),
		http.FxModule(),
		log.FxModule(),
		manager.FxModule(),

		fx.Invoke(http.Run),
	).Run()
}
