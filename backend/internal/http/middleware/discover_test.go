package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDiscoverResource(t *testing.T) {
	// Given
	app := fiber.New()
	app.Use(
		DiscoverResource("my-resource-kind", "my-action"),
		func(c *fiber.Ctx) error {
			ctx := c.UserContext()

			resourceKindKey := ctx.Value(ResourceKindKey).(string)
			resourceValueKey := ctx.Value(ResourceValueKey).(string)
			actionKey := ctx.Value(ActionKey).(string)

			assert.Equal(t, "my-resource-kind", resourceKindKey)
			assert.Equal(t, "*", resourceValueKey)
			assert.Equal(t, "my-action", actionKey)

			return c.Next()
		},
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
