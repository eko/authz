package http

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/http/handler"
)

type Server struct {
	cfg      *configs.HTTPServer
	logger   *slog.Logger
	app      *fiber.App
	handlers handler.Handlers
}

func NewServer(
	app *fiber.App,
	cfg *configs.HTTPServer,
	logger *slog.Logger,
	handlers handler.Handlers,
) *Server {
	server := &Server{
		cfg:      cfg,
		logger:   logger,
		app:      app,
		handlers: handlers,
	}

	server.setSwagger()
	server.setRoutes()

	return server
}

func Run(lc fx.Lifecycle, logger *slog.Logger, server *Server) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting API server", slog.String("addr", server.cfg.Addr))

			go func() {
				if err := server.app.Listen(server.cfg.Addr); err != nil {
					if err != http.ErrServerClosed {
						logger.Error("Unable to serve HTTP API", err)
					}
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("Stopping API server")

			return server.app.Shutdown()
		},
	})
}
