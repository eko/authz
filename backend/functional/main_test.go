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

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/fixtures"
	"github.com/spf13/pflag"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	initializer fixtures.Initializer
	opts        = godog.Options{
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

	app := FxApp()

	app.Start(ctx)

	status := godog.TestSuite{
		Name:                configs.ApplicationName,
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
		authz_clients,
		authz_users,
		authz_oauth_tokens,
		authz_stats,
		authz_actions RESTART IDENTITY CASCADE
		;`).Error; err != nil {
			l.Fatalf("Unable to truncate tables: %v\n", err)
		}

		if err := api.reset(sc); err != nil {
			l.Fatalf("Cannot reset api: %v\n", err)
		}

		if err := initializer.Initialize(); err != nil {
			l.Fatalf("Cannot initialize fixtures: %v\n", err)
		}

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, nil
	})

	ctx.Step(`^I wait "([^"]*)"$`, func(value string) error {
		duration, err := time.ParseDuration(value)
		if err != nil {
			return err
		}

		time.Sleep(duration)
		return nil
	})
	ctx.Step(`^I authenticate with username "([^"]*)" and password "([^"]*)"$`, api.iAuthenticateWithUsernameAndPassword)
	ctx.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, api.iSendRequestTo)
	ctx.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)" with payload:$`, api.iSendRequestToWithPayload)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
}
