package client

import (
	"context"
	"log"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/eko/authz/backend/configs"
	"golang.org/x/oauth2"
)

type Manager interface {
	GetConfig() *oauth2.Config
	GetCookiesDomainName() string
	GetFrontendRedirectURL() string
	GetVerifier() *oidc.IDTokenVerifier
	IsEnabled() bool
}

type manager struct {
	cookiesDomainName   string
	frontendRedirectURL string
	oidcConfig          *oauth2.Config
	oidcProvider        *oidc.Provider
	oidcVerifier        *oidc.IDTokenVerifier
}

func NewManager(
	ctx context.Context,
	cfg *configs.OAuth,
) Manager {
	if cfg.IssuerURL == "" {
		return &manager{}
	}

	oidcProvider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		log.Fatal(err)
	}

	oidcVerifier := oidcProvider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	oidcConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  cfg.RedirectURL,
		Scopes:       append([]string{oidc.ScopeOpenID}, cfg.Scopes...),
	}

	return &manager{
		cookiesDomainName:   cfg.CookiesDomainName,
		frontendRedirectURL: cfg.FrontendRedirectURL,
		oidcConfig:          oidcConfig,
		oidcProvider:        oidcProvider,
		oidcVerifier:        oidcVerifier,
	}
}

func (c *manager) GetConfig() *oauth2.Config {
	return c.oidcConfig
}

func (c *manager) GetCookiesDomainName() string {
	return c.cookiesDomainName
}

func (c *manager) GetFrontendRedirectURL() string {
	return c.frontendRedirectURL
}

func (c *manager) GetVerifier() *oidc.IDTokenVerifier {
	return c.oidcVerifier
}

func (c *manager) IsEnabled() bool {
	return c.oidcProvider != nil
}
