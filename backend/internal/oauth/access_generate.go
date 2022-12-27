package oauth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/security/paseto"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
)

type AccessGenerate struct {
	cfg          *configs.Auth
	tokenManager paseto.Manager
}

func NewAccessGenerate(
	cfg *configs.Auth,
	tokenManager paseto.Manager,
) *AccessGenerate {
	return &AccessGenerate{
		cfg:          cfg,
		tokenManager: tokenManager,
	}
}

func (g *AccessGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	clientID := data.Client.GetID()

	accessToken := g.tokenManager.Generate(clientID, g.cfg.AccessTokenDuration)

	refresh := ""
	if isGenRefresh {
		refresh = base64.URLEncoding.EncodeToString([]byte(uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(clientID)).String()))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return accessToken.Token, refresh, nil
}
