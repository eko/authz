package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/helper/token"
	"github.com/eko/authz/backend/internal/oauth/client"
	"github.com/eko/authz/backend/internal/security/jwt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	OAuthClaimEmailKey = "email"
	OAuthClaimNameKey  = "name"

	OAuthStateCookieName     = "authz_state"
	OAuthExpiresInCookieName = "authz_expires_in"
	OAuthTokenCookieName     = "authz_access_token"
	OAuthNonceCookieName     = "authz_nonce"
)

// Authenticates a user using an OAuth OpenID Connect provider
//
//	@security	Authentication
//	@Summary	Authenticates a user using an OAuth OpenID Connect provider
//	@Tags		Auth
//	@Success	302
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/oauth [Get]
func OAuthAuthenticate(
	oauthClientManager client.Manager,
	tokenGenerator token.Generator,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		state, err := tokenGenerator.Generate(16)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, fmt.Errorf("unable to generate state: %v", err))
		}

		nonce, err := tokenGenerator.Generate(16)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, fmt.Errorf("unable to generate none: %v", err))
		}

		setCallbackCookie(c, OAuthStateCookieName, state, oauthClientManager.GetCookiesDomainName())
		setCallbackCookie(c, OAuthNonceCookieName, nonce, oauthClientManager.GetCookiesDomainName())

		return c.Redirect(
			oauthClientManager.GetConfig().AuthCodeURL(state, oidc.Nonce(nonce)),
			http.StatusFound,
		)
	}
}

// Callback of the OAuth OpenID Connect provider authentication
//
//	@security	Authentication
//	@Summary	Callback of the OAuth OpenID Connect provider authentication
//	@Tags		Auth
//	@Success	200	{object}	AuthResponse
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/oauth/callback [Get]
func OAuthCallback(
	jwtManager jwt.Manager,
	oauthClientManager client.Manager,
	principalManager manager.Principal,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		state := c.Request().Header.Cookie(OAuthStateCookieName)

		if c.Query("state") != string(state) {
			return returnError(c, http.StatusBadRequest, errors.New("state did not match"))
		}

		oauth2Token, err := oauthClientManager.GetConfig().Exchange(ctx, c.Query("code"))
		if err != nil {
			return returnError(c, http.StatusInternalServerError, fmt.Errorf("failed to exchange token: %v", err))
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			return returnError(c, http.StatusInternalServerError, errors.New("no id_token field in oauth2 token"))
		}

		idToken, err := oauthClientManager.GetVerifier().Verify(ctx, rawIDToken)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, fmt.Errorf("failed to verify id token: %v", err))
		}

		nonce := c.Request().Header.Cookie(OAuthNonceCookieName)

		if idToken.Nonce != string(nonce) {
			return returnError(c, http.StatusBadRequest, fmt.Errorf("nonce did not match: %v", err))
		}

		// Obtain user claims from OpenID Connect ID token.
		idTokenClaims := map[string]any{}

		if err := idToken.Claims(&idTokenClaims); err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		emailValue, err := retrieveClaim(idTokenClaims, OAuthClaimEmailKey)
		if err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		nameValue, err := retrieveClaim(idTokenClaims, OAuthClaimNameKey)
		if err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Retrieve or create principal from user email.
		_, err = principalManager.GetRepository().Get(
			model.UserPrincipal(emailValue),
		)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			if _, err = principalManager.Create(
				model.UserPrincipal(emailValue),
				[]string{},
				map[string]any{
					"name": nameValue,
				},
			); err != nil {
				return returnError(c, http.StatusInternalServerError, fmt.Errorf("unable to create principal: %v", err))
			}
		} else if err != nil {
			return returnError(c, http.StatusInternalServerError, fmt.Errorf("unable to retrieve principal: %v", err))
		}

		// Generate access token.
		jwtToken, err := jwtManager.Generate(emailValue)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, fmt.Errorf("unable to generate jwt token: %v", err))
		}

		setCallbackCookie(c, OAuthExpiresInCookieName, strconv.FormatInt(jwtToken.ExpiresIn, 10), oauthClientManager.GetCookiesDomainName())
		setCallbackCookie(c, OAuthTokenCookieName, jwtToken.Token, oauthClientManager.GetCookiesDomainName())

		return c.Redirect(oauthClientManager.GetFrontendRedirectURL(), http.StatusFound)
	}
}

func retrieveClaim(claims map[string]any, key string) (string, error) {
	value, ok := claims[key]
	if !ok {
		return "", fmt.Errorf("unable to retrieve claim from issuer: %s", key)
	}

	return value.(string), nil
}

func setCallbackCookie(c *fiber.Ctx, name, value, domain string) {
	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		HTTPOnly: false,
	})
}
