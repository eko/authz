package middleware

import (
	"log/slog"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/security/jwt"
	"github.com/gofiber/fiber/v3"
)

const (
	AuthenticationKey = "authentication"
	AuthorizationKey  = "authorization"
)

type Middlewares map[string]fiber.Handler

func (m Middlewares) Get(name string) fiber.Handler {
	return m[name]
}

func NewMiddlewares(
	logger *slog.Logger,
	compiledManager manager.CompiledPolicy,
	tokenManager jwt.Manager,
) Middlewares {
	return Middlewares{
		AuthenticationKey: Authentication(logger, tokenManager),
		AuthorizationKey:  Authorization(logger, compiledManager),
	}
}
