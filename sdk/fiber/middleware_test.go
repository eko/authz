package fiber

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eko/authz/backend/pkg/authz"
	gofiber "github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type checkerFunc func(ctx context.Context, check *authz.Check) (bool, error)

func (f checkerFunc) IsAllowed(ctx context.Context, check *authz.Check) (bool, error) {
	return f(ctx, check)
}

func TestNew(t *testing.T) {
	principalFromHeader := func(c gofiber.Ctx) (string, error) {
		return c.Get("X-User"), nil
	}

	testCases := []struct {
		name          string
		user          string
		path          string
		allowed       bool
		checkErr      error
		expectedCode  int
		expectedCheck *authz.Check
	}{
		{
			name:         "allowed on a specific resource",
			user:         "user-123",
			path:         "/posts/456",
			allowed:      true,
			expectedCode: http.StatusOK,
			expectedCheck: &authz.Check{
				Principal: "user-123", ResourceKind: "post", ResourceValue: "456", Action: "edit",
			},
		},
		{
			name:         "allowed on wildcard resource",
			user:         "user-123",
			path:         "/posts",
			allowed:      true,
			expectedCode: http.StatusOK,
			expectedCheck: &authz.Check{
				Principal: "user-123", ResourceKind: "post", ResourceValue: "*", Action: "edit",
			},
		},
		{
			name:         "denied",
			user:         "user-123",
			path:         "/posts/456",
			allowed:      false,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "unknown principal",
			user:         "",
			path:         "/posts/456",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "check error fails closed",
			user:         "user-123",
			path:         "/posts/456",
			allowed:      true,
			checkErr:     errors.New("connection refused"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var received *authz.Check

			middleware := New(Config{
				Checker: checkerFunc(func(_ context.Context, check *authz.Check) (bool, error) {
					received = check
					return testCase.allowed, testCase.checkErr
				}),
				Principal:     principalFromHeader,
				ResourceKind:  "post",
				ResourceValue: ResourceValueFromParam("id"),
				Action:        "edit",
			})

			ok := func(c gofiber.Ctx) error { return c.SendStatus(http.StatusOK) }

			app := gofiber.New()
			app.Get("/posts", middleware, ok)
			app.Get("/posts/:id", middleware, ok)

			req := httptest.NewRequest(http.MethodGet, testCase.path, nil)
			req.Header.Set("X-User", testCase.user)

			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedCode, resp.StatusCode)
			if testCase.expectedCheck != nil {
				assert.Equal(t, testCase.expectedCheck.String(), received.String())
			}
		})
	}
}

func TestNew_PanicsOnInvalidConfig(t *testing.T) {
	assert.Panics(t, func() { New(Config{}) })
}
