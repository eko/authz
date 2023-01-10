package sdk

import (
	"context"
	"fmt"

	"github.com/eko/authz/backend/pkg/authz"
	"github.com/eko/authz/sdk/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	authz.ApiClient
	IsAllowed(ctx context.Context, check *authz.Check) (bool, error)
}

type client struct {
	cfg *Config
	authz.ApiClient
}

func NewClient(cfg *Config) (Client, error) {
	if cfg == nil {
		cfg = DefaultConfig
	}

	authenticator := interceptor.NewAuthenticator(
		cfg.ClientID, cfg.ClientSecret,
		interceptor.WithExpireDelay(cfg.ExpireDelay),
	)

	clientConn, err := grpc.Dial(cfg.GrpcAddr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithChainUnaryInterceptor(
			interceptor.AuthenticationUnaryClientInterceptor(authenticator),
		),
		grpc.WithChainStreamInterceptor(
			interceptor.AuthenticationStreamClientInterceptor(authenticator),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to gRPC backend: %v", err)
	}

	apiClient := authz.NewApiClient(clientConn)

	return &client{
		cfg:       cfg,
		ApiClient: apiClient,
	}, nil
}

func (c *client) IsAllowed(ctx context.Context, check *authz.Check) (bool, error) {
	if check == nil {
		return false, nil
	}

	response, err := c.Check(ctx, &authz.CheckRequest{
		Checks: []*authz.Check{check},
	})
	if err != nil {
		return false, err
	}

	return response.GetChecks()[0].IsAllowed, nil
}
