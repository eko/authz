package interceptor

import (
	"context"
	"slices"
	"time"

	"github.com/eko/authz/backend/pkg/authz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationKey = "authorization"
	authenticateRPC  = "/authz.Api/Authenticate"
)

var (
	// UnauthenticatedMethods specify gRPC methods that should not be authenticated.
	// This means they can be call publicly.
	UnauthenticatedMethods = []string{
		authenticateRPC,
		"/authz.Api/Check",
	}
)

// accessToken returns a valid access token, authenticating against the
// backend only when no token is cached or when it is about to expire.
func accessToken(ctx context.Context, cc *grpc.ClientConn, authenticator Authenticator, opts ...grpc.CallOption) (string, error) {
	if token := authenticator.GetToken(); token.IsValid(time.Now().Add(authenticator.GetExpireDelay())) {
		return token.AccessToken, nil
	}

	request := &authz.AuthenticateRequest{
		ClientId:     authenticator.GetClientID(),
		ClientSecret: authenticator.GetClientSecret(),
	}
	response := &authz.AuthenticateResponse{}
	if err := cc.Invoke(ctx, authenticateRPC, request, response, opts...); err != nil {
		return "", err
	}

	token := &Token{
		AccessToken: response.GetToken(),
		ExpireAt:    time.Now().Add(time.Duration(response.GetExpiresIn()) * time.Second),
	}
	authenticator.SetToken(token)

	return token.AccessToken, nil
}

func authenticatedContext(ctx context.Context, cc *grpc.ClientConn, authenticator Authenticator, opts ...grpc.CallOption) (context.Context, error) {
	token, err := accessToken(ctx, cc, authenticator, opts...)
	if err != nil {
		return ctx, err
	}

	return metadata.AppendToOutgoingContext(ctx, authorizationKey, "bearer "+token), nil
}

func AuthenticationUnaryClientInterceptor(authenticator Authenticator) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if slices.Contains(UnauthenticatedMethods, method) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		newCtx, err := authenticatedContext(ctx, cc, authenticator, opts...)
		if err != nil {
			return err
		}

		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

func AuthenticationStreamClientInterceptor(authenticator Authenticator) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if slices.Contains(UnauthenticatedMethods, method) {
			return streamer(ctx, desc, cc, method, opts...)
		}

		newCtx, err := authenticatedContext(ctx, cc, authenticator, opts...)
		if err != nil {
			return nil, err
		}

		return streamer(newCtx, desc, cc, method, opts...)
	}
}
