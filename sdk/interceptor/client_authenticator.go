package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/eko/authz/backend/pkg/authz"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
)

const (
	authorizationKey = "authorization"
)

var (
	// UnauthenticatedMethods specify gRPC methods that should not be authenticated.
	// This means they can be call publicly.
	UnauthenticatedMethods = []string{
		"/authz.Api/Authenticate",
		"/authz.Api/Check",
	}
)

func checkAndApply(ctx context.Context, cc *grpc.ClientConn, authenticator Authenticator, opts ...grpc.CallOption) (context.Context, error) {
	authenticateRequest := &authz.AuthenticateRequest{
		ClientId:     authenticator.GetClientID(),
		ClientSecret: authenticator.GetClientSecret(),
	}
	authenticateResponse := &authz.AuthenticateResponse{}
	if err := cc.Invoke(ctx, "/authz.Api/Authenticate", authenticateRequest, authenticateResponse, opts...); err != nil {
		return ctx, err
	}

	if authenticator.GetToken() == nil ||
		authenticator.GetToken().ExpireAt.IsZero() ||
		authenticator.GetToken().ExpireAt.Before(time.Now().Add(authenticator.GetExpireDelay())) {

		authenticator.SetToken(&Token{
			AccessToken: authenticateResponse.GetToken(),
			ExpireAt:    time.Now().Add(time.Duration(authenticateResponse.GetExpiresIn()) * time.Second),
		})
	}

	newCtx := metautils.
		ExtractIncoming(ctx).
		Set(authorizationKey, fmt.Sprintf("bearer %s", authenticator.GetToken().AccessToken)).
		ToOutgoing(ctx)

	return newCtx, nil
}

func AuthenticationUnaryClientInterceptor(authenticator Authenticator) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if slices.Contains(UnauthenticatedMethods, method) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		newCtx, err := checkAndApply(ctx, cc, authenticator, opts...)
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

		newCtx, err := checkAndApply(ctx, cc, authenticator, opts...)
		if err != nil {
			return nil, err
		}

		return streamer(newCtx, desc, cc, method, opts...)
	}
}
