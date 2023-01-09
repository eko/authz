package handler

import (
	"context"
	"errors"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/security/jwt"
	"github.com/eko/authz/backend/pkg/authz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	// ClientNotFoundErr is returned when client_id does not exists.
	ClientNotFoundErr = status.Error(codes.NotFound, "client not found")

	// InvalidCredentialsErr is returned when client_id or client_secret is invalid.
	InvalidCredentialsErr = status.Error(codes.InvalidArgument, "invalid credentials")
)

type Auth interface {
	Authenticate(ctx context.Context, req *authz.AuthenticateRequest) (*authz.AuthenticateResponse, error)
}

type auth struct {
	clientManager manager.Client
	tokenManager  jwt.Manager
}

func NewAuth(
	clientManager manager.Client,
	tokenManager jwt.Manager,
) Auth {
	return &auth{
		clientManager: clientManager,
		tokenManager:  tokenManager,
	}
}

func (h *auth) Authenticate(ctx context.Context, req *authz.AuthenticateRequest) (*authz.AuthenticateResponse, error) {
	client, err := h.clientManager.GetRepository().Get(req.GetClientId())
	if err != nil || client.Secret != req.GetClientSecret() {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ClientNotFoundErr
		}

		return nil, InvalidCredentialsErr
	}

	token, err := h.tokenManager.Generate(client.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authz.AuthenticateResponse{
		Token:     token.Token,
		Type:      token.TokenType,
		ExpiresIn: token.ExpiresIn,
	}, nil
}
