package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type contextKey string

var (
	ResourceKindKey  = contextKey("authz_resource_kind")
	ResourceValueKey = contextKey("authz_resource_value")
	ActionKey        = contextKey("authz_action")
)

func DiscoverResource(resourceKind string, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		ctx = context.WithValue(ctx, ResourceKindKey, resourceKind)
		ctx = context.WithValue(ctx, ResourceValueKey, c.Params("identifier", "*"))
		ctx = context.WithValue(ctx, ActionKey, action)

		c.SetUserContext(ctx)

		return c.Next()
	}
}
