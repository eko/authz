package middleware

import (
	"fmt"

	"github.com/eko/authz/backend/internal/manager"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

func Authorization(
	logger *slog.Logger,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		userID := ctx.Value(UserIdentifierKey).(string)
		resourceKind := ctx.Value(ResourceKindKey).(string)
		resourceValue := ctx.Value(ResourceValueKey).(string)
		action := ctx.Value(ActionKey).(string)

		principal := fmt.Sprintf("authz-%s", userID)

		isAllowed, err := manager.IsAllowed(principal, resourceKind, resourceValue, action)
		if err != nil {
			logger.Error(
				"Error while checking if user is allowed",
				err,
				slog.String("principal", principal),
				slog.String("resource_kind", resourceKind),
				slog.String("resource_value", resourceValue),
				slog.String("action", action),
			)
		}

		if !isAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "access denied",
			})
		}

		return c.Next()
	}
}
