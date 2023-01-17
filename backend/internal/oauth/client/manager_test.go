package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/eko/authz/backend/configs"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var (
	mockedOpenIDConfigurationHandler = func(issuerUrl string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			_, _ = w.Write([]byte(`{
			"issuer": "` + issuerUrl + `"
		}`))
		}
	}
)

func TestNewManager(t *testing.T) {
	// Given
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Config.Handler = mockedOpenIDConfigurationHandler(server.URL)

	cfg := &configs.OAuth{
		IssuerURL: server.URL,
	}

	// When
	manager := NewManager(ctx, cfg)

	// Then
	assert.Implements(t, new(Manager), manager)
}

func TestManager_GetConfig(t *testing.T) {
	// Given
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Config.Handler = mockedOpenIDConfigurationHandler(server.URL)

	cfg := &configs.OAuth{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		IssuerURL:    server.URL,
		RedirectURL:  "http://localhost:8080/v1/oauth/callback",
		Scopes:       []string{"test-scope-1", "test-scope-2"},
	}

	manager := NewManager(ctx, cfg)

	// When
	config := manager.GetConfig()

	// Then
	assert.IsType(t, new(oauth2.Config), config)

	assert.Equal(t, cfg.ClientID, config.ClientID)
	assert.Equal(t, cfg.ClientSecret, config.ClientSecret)
	assert.Equal(t, cfg.RedirectURL, config.RedirectURL)
	assert.Equal(t, []string{
		oidc.ScopeOpenID,
		"test-scope-1",
		"test-scope-2",
	}, config.Scopes)
}

func TestManager_GetCookiesDomainName(t *testing.T) {
	// Given
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Config.Handler = mockedOpenIDConfigurationHandler(server.URL)

	cfg := &configs.OAuth{
		CookiesDomainName: "test-cookies-domain-name.acme.tld",
		IssuerURL:         server.URL,
	}

	manager := NewManager(ctx, cfg)

	// When - Then
	assert.Equal(t, cfg.CookiesDomainName, manager.GetCookiesDomainName())
}

func TestManager_GetFrontendRedirectURL(t *testing.T) {
	// Given
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Config.Handler = mockedOpenIDConfigurationHandler(server.URL)

	cfg := &configs.OAuth{
		FrontendRedirectURL: "http://test.acme.tld",
		IssuerURL:           server.URL,
	}

	manager := NewManager(ctx, cfg)

	// When - Then
	assert.Equal(t, cfg.FrontendRedirectURL, manager.GetFrontendRedirectURL())
}

func TestManager_GetVerifier(t *testing.T) {
	// Given
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Config.Handler = mockedOpenIDConfigurationHandler(server.URL)

	cfg := &configs.OAuth{
		ClientID:  "test-client-id",
		IssuerURL: server.URL,
	}

	manager := NewManager(ctx, cfg)

	// When
	verifier := manager.GetVerifier()

	// Then
	assert.IsType(t, new(oidc.IDTokenVerifier), verifier)
}

func TestManager_IsEnabled_True(t *testing.T) {
	// Given
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Config.Handler = mockedOpenIDConfigurationHandler(server.URL)

	cfg := &configs.OAuth{
		IssuerURL: server.URL,
	}

	manager := NewManager(ctx, cfg)

	// When - Then
	assert.True(t, manager.IsEnabled())
}

func TestManager_IsEnabled_False(t *testing.T) {
	// Given
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Config.Handler = mockedOpenIDConfigurationHandler(server.URL)

	cfg := &configs.OAuth{}

	manager := NewManager(ctx, cfg)

	// When - Then
	assert.False(t, manager.IsEnabled())
}
