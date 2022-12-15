//go:build functional
// +build functional

package main

import (
	"context"
	l "log"
	lib_http "net/http"
	"os"
	"testing"
	"time"

	_ "time/tzdata"

	"bou.ke/monkey"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/compile"
	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/helper"
	"github.com/eko/authz/backend/internal/http"
	"github.com/eko/authz/backend/internal/log"
	"github.com/eko/authz/backend/internal/manager"
	"github.com/spf13/pflag"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	opts = godog.Options{
		Output: colors.Colored(os.Stdout),
	}
)

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	pflag.Parse()
	opts.Paths = pflag.Args()

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		l.Fatalf("Unable to load location: %v\n", err)
	}
	time.Local = loc

	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2100, time.January, 1, 1, 0, 0, 0, time.UTC)
	})

	app := fx.New(
		fx.NopLogger,
		fx.Provide(context.Background),

		compile.FxModule(),
		database.FxModule(),
		event.FxModule(),
		helper.FxModule(),
		http.FxModule(),
		log.FxModule(),
		manager.FxModule(),

		fx.Provide(
			configs.Load,
			func(cfg *configs.Base) *configs.Database { return cfg.Database },
			func(cfg *configs.Base) *configs.HTTPServer { return cfg.HTTPServer },
			func(cfg *configs.Base) *configs.Logger {
				cfg.Logger.Level = "ERROR"
				return cfg.Logger
			},
		),

		fx.Invoke(
			func(database *gorm.DB) { db = database },
		),

		fx.Invoke(http.Run),
	)

	app.Start(ctx)

	status := godog.TestSuite{
		Name:                "authz",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if status != 0 {
		os.Exit(1)
	}

	_ = app.Stop(ctx)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{
		httpClient: &lib_http.Client{},
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		if err := db.Exec(`TRUNCATE TABLE
		authz_compiled_policies,
		authz_roles_policies,
		authz_roles,
		authz_principals_roles,
		authz_principals,
		authz_policies_actions,
		authz_policies_resources,
		authz_policies,
		authz_resources,
		authz_actions RESTART IDENTITY CASCADE
		;`).Error; err != nil {
			l.Fatalf("Unable to truncate tables: %v\n", err)
		}

		if err := api.reset(sc); err != nil {
			l.Fatalf("Cannot reset api: %v\n", err)
		}

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, nil
	})

	ctx.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, api.iSendRequestTo)
	ctx.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)" with payload:$`, api.iSendRequestToWithPayload)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
}
