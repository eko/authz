package interceptor

import (
	"context"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/security/jwt"
)

type AuthzFunc func(ctx context.Context, resourceKind string, resourceValue string, action string) bool

func AuthorizationFunc(compiledManager manager.CompiledPolicy) AuthzFunc {
	return func(ctx context.Context, resourceKind string, resourceValue string, action string) bool {
		claims, ok := ctx.Value(ClaimsKey).(*jwt.Claims)
		if !ok {
			return false
		}

		isAllowed, err := compiledManager.IsAllowed(claims.Subject, resourceKind, resourceValue, action)
		if err != nil {
			return false
		}

		return isAllowed
	}
}
