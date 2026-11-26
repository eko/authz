package middleware

import (
	"context"
	"log/slog"
	"strings"

	"github.com/eko/authz/backend/internal/security/jwt"
	"github.com/gofiber/fiber/v3"
)

var (
	UserIdentifierKey = struct{}{}
)

func Authentication(
	logger *slog.Logger,
	tokenManager jwt.Manager,
) fiber.Handler {
	return func(c fiber.Ctx) error {
		authorizationHeader := c.Get("Authorization")
		authorizationHeaderValues := strings.Split(authorizationHeader, " ")

		if authorizationHeader == "" || len(authorizationHeaderValues) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "unauthorized",
			})
		}

		token := authorizationHeaderValues[1]

		claims, err := tokenManager.Parse(token)
		if err != nil {
			logger.Error("unable to verify token", slog.Any("error", err))

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "unable to verify token",
			})
		}

		ctx := c.Context()
		ctx = context.WithValue(ctx, UserIdentifierKey, claims.Subject)

		c.SetContext(ctx)

		return c.Next()
	}
}
