package middleware

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eko/authz/backend/internal/log"
	"github.com/eko/authz/backend/internal/security/jwt"
	"github.com/gofiber/fiber/v2"
	lib_jwt "github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestAuthentication_WhenNoAuthorizationHeader(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	jwtManager := jwt.NewMockManager(ctrl)

	app := fiber.New()
	app.Use(Authentication(logger, jwtManager))

	// When
	req := httptest.NewRequest("GET", "/", nil)

	response, err := app.Test(req)
	assert.Nil(t, err)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, `{"error":true,"message":"unauthorized"}`, string(bodyBytes))
}

func TestAuthentication_WhenMalformedAuthorizationHeader(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	jwtManager := jwt.NewMockManager(ctrl)

	app := fiber.New()
	app.Use(Authentication(logger, jwtManager))

	// When
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Something that is not a bearer token")

	response, err := app.Test(req)
	assert.Nil(t, err)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, `{"error":true,"message":"unauthorized"}`, string(bodyBytes))
}

func TestAuthentication_WhenTokenParseError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	expectedErr := errors.New("token error")

	jwtManager := jwt.NewMockManager(ctrl)
	jwtManager.EXPECT().Parse("token-123").Return(nil, expectedErr)

	app := fiber.New()
	app.Use(Authentication(logger, jwtManager))

	// When
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer token-123")

	response, err := app.Test(req)
	assert.Nil(t, err)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, `{"error":true,"message":"unable to verify token"}`, string(bodyBytes))
}

func TestAuthentication_WhenValidToken(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	jwtManager := jwt.NewMockManager(ctrl)
	jwtManager.EXPECT().Parse("token-123").Return(&jwt.Claims{
		RegisteredClaims: lib_jwt.RegisteredClaims{
			Subject: "user-123",
		},
	}, nil)

	app := fiber.New()
	app.Use(Authentication(logger, jwtManager))
	app.Get("/", func(c *fiber.Ctx) error {
		_ = c.JSON(map[string]any{"success": true})
		return nil
	})

	// When
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer token-123")

	response, err := app.Test(req)
	assert.Nil(t, err)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	// Then
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, `{"success":true}`, string(bodyBytes))
}
