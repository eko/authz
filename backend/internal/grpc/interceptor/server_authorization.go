package interceptor

import (
	"context"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/pkg/authz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ResourcesAndActionsByMethod maps the resource kind and action for each
	// gRPC method available in the proto API.
	ResourcesAndActionsByMethod = map[string][]string{
		"/authz.Api/PolicyCreate": {"authz.policies", "create"},
		"/authz.Api/PolicyDelete": {"authz.policies", "delete"},
		"/authz.Api/PolicyGet":    {"authz.policies", "get"},
		"/authz.Api/PolicyUpdate": {"authz.policies", "update"},

		"/authz.Api/PrincipalCreate": {"authz.principals", "create"},
		"/authz.Api/PrincipalDelete": {"authz.principals", "delete"},
		"/authz.Api/PrincipalGet":    {"authz.principals", "get"},
		"/authz.Api/PrincipalUpdate": {"authz.principals", "update"},

		"/authz.Api/ResourceCreate": {"authz.resources", "create"},
		"/authz.Api/ResourceDelete": {"authz.resources", "delete"},
		"/authz.Api/ResourceGet":    {"authz.resources", "get"},
		"/authz.Api/ResourceUpdate": {"authz.resources", "update"},

		"/authz.Api/RoleCreate": {"authz.roles", "create"},
		"/authz.Api/RoleDelete": {"authz.roles", "delete"},
		"/authz.Api/RoleGet":    {"authz.roles", "get"},
		"/authz.Api/RoleUpdate": {"authz.roles", "update"},
	}

	// RetrieveResourceValueByMethod maps the request object for each gRPC method
	// that needs a resource value (identifier).
	RetrieveResourceValueByMethod = map[string]string{
		"/authz.Api/PolicyDelete": "PolicyDeleteRequest",
		"/authz.Api/PolicyGet":    "PolicyGetRequest",
		"/authz.Api/PolicyUpdate": "PolicyUpdateRequest",

		"/authz.Api/PrincipalDelete": "PrincipalDeleteRequest",
		"/authz.Api/PrincipalGet":    "PrincipalGetRequest",
		"/authz.Api/PrincipalUpdate": "PrincipalUpdateRequest",

		"/authz.Api/ResourceDelete": "ResourceDeleteRequest",
		"/authz.Api/ResourceGet":    "ResourceGetRequest",
		"/authz.Api/ResourceUpdate": "ResourceUpdateRequest",

		"/authz.Api/RoleDelete": "RoleDeleteRequest",
		"/authz.Api/RoleGet":    "RoleGetRequest",
		"/authz.Api/RoleUpdate": "RoleUpdateRequest",
	}
)

// AuthorizationStreamServerInterceptor checks if current user is allowed to do stream method calls.
func AuthorizationStreamServerInterceptor(authorizationFunc AuthzFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		values, ok := ResourcesAndActionsByMethod[info.FullMethod]
		if !ok {
			// Method authorization is not managed, let client continue.
			return handler(srv, stream)
		}

		resourceKind, action := values[0], values[1]

		if !authorizationFunc(stream.Context(), resourceKind, manager.WildcardValue, action) {
			return status.Errorf(codes.PermissionDenied, "not allowed to do %s on %s", action, resourceKind)
		}

		return handler(srv, stream)
	}
}

// AuthorizationUnaryServerInterceptor checks if current user is allowed to do method calls.
func AuthorizationUnaryServerInterceptor(authorizationFunc AuthzFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		values, ok := ResourcesAndActionsByMethod[info.FullMethod]
		if !ok {
			// Method authorization is not managed, let client continue.
			return handler(ctx, req)
		}

		resourceKind, action := values[0], values[1]
		resourceValue := getResourceValue(info.FullMethod, req)

		if !authorizationFunc(ctx, resourceKind, resourceValue, action) {
			return nil, status.Errorf(codes.PermissionDenied, "not allowed to do %s on %s (%s)", action, resourceKind, resourceValue)
		}

		return handler(ctx, req)
	}
}

func getResourceValue(method string, req any) string {
	value, ok := RetrieveResourceValueByMethod[method]
	if !ok {
		return manager.WildcardValue
	}

	switch value {
	case "PolicyDeleteRequest":
		return req.(*authz.PolicyDeleteRequest).GetId()
	case "PolicyGetRequest":
		return req.(*authz.PolicyGetRequest).GetId()
	case "PolicyUpdateRequest":
		return req.(*authz.PolicyUpdateRequest).GetId()

	case "PrincipalDeleteRequest":
		return req.(*authz.PrincipalDeleteRequest).GetId()
	case "PrincipalGetRequest":
		return req.(*authz.PrincipalGetRequest).GetId()
	case "PrincipalUpdateRequest":
		return req.(*authz.PrincipalUpdateRequest).GetId()

	case "ResourceDeleteRequest":
		return req.(*authz.ResourceDeleteRequest).GetId()
	case "ResourceGetRequest":
		return req.(*authz.ResourceGetRequest).GetId()
	case "ResourceUpdateRequest":
		return req.(*authz.ResourceUpdateRequest).GetId()

	case "RoleDeleteRequest":
		return req.(*authz.RoleDeleteRequest).GetId()
	case "RoleGetRequest":
		return req.(*authz.RoleGetRequest).GetId()
	case "RoleUpdateRequest":
		return req.(*authz.RoleUpdateRequest).GetId()
	}

	return manager.WildcardValue
}
