package oauth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/eko/authz/backend/internal/security/jwt"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
)

type AccessGenerate struct {
	tokenManager jwt.Manager
}

func NewAccessGenerate(
	tokenManager jwt.Manager,
) *AccessGenerate {
	return &AccessGenerate{
		tokenManager: tokenManager,
	}
}

func (g *AccessGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	clientID := data.Client.GetID()

	accessToken, err := g.tokenManager.Generate(clientID)
	if err != nil {
		return "", "", err
	}

	refresh := ""
	if isGenRefresh {
		refresh = base64.URLEncoding.EncodeToString([]byte(uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(clientID)).String()))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return accessToken.Token, refresh, nil
}
