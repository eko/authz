package http

import (
	"github.com/eko/authz/backend/internal/http/handler"
	handler_validator "github.com/eko/authz/backend/internal/http/handler/validator"
	"github.com/eko/authz/backend/internal/http/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("http",
		fx.Provide(
			func() *fiber.App {
				app := fiber.New()
				// A panic in a handler must not bring the whole server down.
				app.Use(recover.New())

				return app
			},
			NewServer,
			handler.NewHandlers,
			middleware.NewMiddlewares,
			validator.New,
		),
		fx.Invoke(handler_validator.Register),
	)
}
