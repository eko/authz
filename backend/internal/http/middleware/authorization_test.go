package middleware

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/log"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestAuthorization_WhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	expectedErr := errors.New("some error")

	compiledManager := manager.NewMockCompiledPolicy(ctrl)
	compiledManager.EXPECT().IsAllowed(
		"authz-user-user-123",
		"my-resource-kind",
		"my-resource-value",
		"my-action",
	).Return(false, expectedErr)

	app := fiber.New()
	app.Use(
		func(c *fiber.Ctx) error {
			ctx := context.Background()
			ctx = context.WithValue(ctx, UserIdentifierKey, "user-123")
			ctx = context.WithValue(ctx, ResourceKindKey, "my-resource-kind")
			ctx = context.WithValue(ctx, ResourceValueKey, "my-resource-value")
			ctx = context.WithValue(ctx, ActionKey, "my-action")

			c.SetUserContext(ctx)
			return c.Next()
		},
		Authorization(logger, compiledManager),
	)
	app.Get("/", func(c *fiber.Ctx) error {
		_ = c.JSON(map[string]any{"success": true})
		return nil
	})

	// When
	req := httptest.NewRequest("GET", "/", nil)

	response, err := app.Test(req)
	assert.Nil(t, err)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
	assert.Equal(t, `{"error":true,"message":"access denied"}`, string(bodyBytes))
}

func TestAuthorization_WhenIsNotAllowed(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	compiledManager := manager.NewMockCompiledPolicy(ctrl)
	compiledManager.EXPECT().IsAllowed(
		"authz-user-user-123",
		"my-resource-kind",
		"my-resource-value",
		"my-action",
	).Return(false, nil)

	app := fiber.New()
	app.Use(
		func(c *fiber.Ctx) error {
			ctx := context.Background()
			ctx = context.WithValue(ctx, UserIdentifierKey, "user-123")
			ctx = context.WithValue(ctx, ResourceKindKey, "my-resource-kind")
			ctx = context.WithValue(ctx, ResourceValueKey, "my-resource-value")
			ctx = context.WithValue(ctx, ActionKey, "my-action")

			c.SetUserContext(ctx)
			return c.Next()
		},
		Authorization(logger, compiledManager),
	)
	app.Get("/", func(c *fiber.Ctx) error {
		_ = c.JSON(map[string]any{"success": true})
		return nil
	})

	// When
	req := httptest.NewRequest("GET", "/", nil)

	response, err := app.Test(req)
	assert.Nil(t, err)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
	assert.Equal(t, `{"error":true,"message":"access denied"}`, string(bodyBytes))
}

func TestAuthorization_WhenIsAllowed(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	compiledManager := manager.NewMockCompiledPolicy(ctrl)
	compiledManager.EXPECT().IsAllowed(
		"authz-user-user-123",
		"my-resource-kind",
		"my-resource-value",
		"my-action",
	).Return(true, nil)

	app := fiber.New()
	app.Use(
		func(c *fiber.Ctx) error {
			ctx := context.Background()
			ctx = context.WithValue(ctx, UserIdentifierKey, "user-123")
			ctx = context.WithValue(ctx, ResourceKindKey, "my-resource-kind")
			ctx = context.WithValue(ctx, ResourceValueKey, "my-resource-value")
			ctx = context.WithValue(ctx, ActionKey, "my-action")

			c.SetUserContext(ctx)
			return c.Next()
		},
		Authorization(logger, compiledManager),
	)
	app.Get("/", func(c *fiber.Ctx) error {
		_ = c.JSON(map[string]any{"success": true})
		return nil
	})

	// When
	req := httptest.NewRequest("GET", "/", nil)

	response, err := app.Test(req)
	assert.Nil(t, err)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, `{"success":true}`, string(bodyBytes))
}
