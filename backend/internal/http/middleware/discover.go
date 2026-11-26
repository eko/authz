package middleware

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

type contextKey string

var (
	ResourceKindKey  = contextKey("authz_resource_kind")
	ResourceValueKey = contextKey("authz_resource_value")
	ActionKey        = contextKey("authz_action")
)

func DiscoverResource(resourceKind string, action string) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.Context()

		ctx = context.WithValue(ctx, ResourceKindKey, resourceKind)
		ctx = context.WithValue(ctx, ResourceValueKey, c.Params("identifier", "*"))
		ctx = context.WithValue(ctx, ActionKey, action)

		c.SetContext(ctx)

		return c.Next()
	}
}
