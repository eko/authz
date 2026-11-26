package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

type unconfiguredOAuthManager struct{}

func (unconfiguredOAuthManager) GetConfig() *oauth2.Config          { return nil }
func (unconfiguredOAuthManager) GetCookiesDomainName() string       { return "" }
func (unconfiguredOAuthManager) GetFrontendRedirectURL() string     { return "" }
func (unconfiguredOAuthManager) GetVerifier() *oidc.IDTokenVerifier { return nil }
func (unconfiguredOAuthManager) IsEnabled() bool                    { return false }

type staticTokenGenerator struct{}

func (staticTokenGenerator) Generate(int) (string, error) { return "token", nil }

func TestOAuthHandlers_WhenOAuthIsNotConfigured(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		handler fiber.Handler
	}{
		{
			name:    "authenticate",
			path:    "/v1/oauth",
			handler: OAuthAuthenticate(unconfiguredOAuthManager{}, staticTokenGenerator{}),
		},
		{
			name:    "callback",
			path:    "/v1/oauth/callback?state=&code=abc",
			handler: OAuthCallback(nil, unconfiguredOAuthManager{}, nil),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/v1/oauth", testCase.handler)
			app.Get("/v1/oauth/callback", testCase.handler)

			resp, err := app.Test(httptest.NewRequest(http.MethodGet, testCase.path, nil))
			require.NoError(t, err)

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	}
}
