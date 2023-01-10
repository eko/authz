package interceptor

import (
	"context"

	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
)

var (
	// UnauthenticatedMethods specify gRPC methods that should not be authenticated.
	// This means they can be call publicly.
	UnauthenticatedMethods = []string{
		"/authz.Api/Authenticate",
		"/authz.Api/Check",
	}
)

// AuthenticationUnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func AuthenticationUnaryServerInterceptor(interceptor grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if slices.Contains(UnauthenticatedMethods, info.FullMethod) {
			return handler(ctx, req)
		}

		return interceptor(ctx, req, info, handler)
	}
}
