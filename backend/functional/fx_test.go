//go:build functional
// +build functional

//nolint:typecheck
package main

import (
	"context"
	l "log"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/compile"
	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/entity"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/fixtures"
	"github.com/eko/authz/backend/internal/helper/time"
	"github.com/eko/authz/backend/internal/helper/token"
	"github.com/eko/authz/backend/internal/http"
	"github.com/eko/authz/backend/internal/log"
	"github.com/eko/authz/backend/internal/oauth"
	"github.com/eko/authz/backend/internal/security"
	"github.com/eko/authz/backend/internal/stats"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

func FxApp() *fx.App {
	return fx.New(
		fx.NopLogger,
		fx.Provide(context.Background),

		compile.FxModule(),
		database.FxModule(),
		entity.FxModule(),
		event.FxModule(),
		http.FxModule(),
		log.FxModule(),
		oauth.FxModule(),
		security.FxModule(),
		stats.FxModule(),

		fx.Provide(func(
			cfg *configs.User,
			logger *slog.Logger,
			policyManager manager.Policy,
			principalManager manager.Principal,
			resourceManager manager.Resource,
			roleManager manager.Role,
			userManager manager.User,
		) fixtures.Initializer {
			initializer = fixtures.NewInitializer(cfg, logger, policyManager, principalManager, resourceManager, roleManager, userManager)
			return initializer
		}),

		fx.Provide(
			time.NewStaticClock,
			func(clock *time.StaticClock) time.Clock { return clock },
			token.NewGenerator,
		),

		fx.Provide(
			configs.Load,
			func(cfg *configs.Base) *configs.App {
				return cfg.App
			},
			func(cfg *configs.Base) *configs.Auth { return cfg.Auth },
			func(cfg *configs.Base) *configs.Database { return cfg.Database },
			func(cfg *configs.Base) *configs.HTTPServer { return cfg.HTTPServer },
			func(cfg *configs.Base) *configs.Logger {
				cfg.Logger.Level = "ERROR"
				return cfg.Logger
			},
			func(cfg *configs.Base) *configs.OAuth { return cfg.OAuth },
			func(cfg *configs.Base) *configs.User { return cfg.User },
		),

		fx.Invoke(
			func(database *gorm.DB) { db = database },
		),

		fx.Invoke(http.Run),
		fx.Invoke(func(initializer fixtures.Initializer) {
			if err := initializer.Initialize(); err != nil {
				l.Fatalf("Cannot initialize fixtures: %v\n", err)
			}
		}),
	)
}
