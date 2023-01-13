package interceptor

import (
	"context"
	"errors"
	"testing"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/security/jwt"
	lib_jwt "github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizationFunc_WhenIsAllowed(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	ctx = context.WithValue(ctx, ClaimsKey, &jwt.Claims{
		RegisteredClaims: lib_jwt.RegisteredClaims{
			Subject: "my-subject",
		},
	})

	resourceKind := "my-resource-kind"
	resourceValue := "my-resource-value"
	action := "my-action"

	compiledManager := manager.NewMockCompiledPolicy(ctrl)
	compiledManager.EXPECT().IsAllowed("my-subject", resourceKind, resourceValue, action).
		Return(true, nil)

	// When
	result := AuthorizationFunc(compiledManager)(ctx, resourceKind, resourceValue, action)

	// Then
	assert.Equal(t, true, result)
}

func TestAuthorizationFunc_WhenIsNotAllowed(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	ctx = context.WithValue(ctx, ClaimsKey, &jwt.Claims{
		RegisteredClaims: lib_jwt.RegisteredClaims{
			Subject: "my-subject",
		},
	})

	resourceKind := "my-resource-kind"
	resourceValue := "my-resource-value"
	action := "my-action"

	compiledManager := manager.NewMockCompiledPolicy(ctrl)
	compiledManager.EXPECT().IsAllowed("my-subject", resourceKind, resourceValue, action).
		Return(false, nil)

	// When
	result := AuthorizationFunc(compiledManager)(ctx, resourceKind, resourceValue, action)

	// Then
	assert.Equal(t, false, result)
}

func TestAuthorizationFunc_WhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	ctx = context.WithValue(ctx, ClaimsKey, &jwt.Claims{
		RegisteredClaims: lib_jwt.RegisteredClaims{
			Subject: "my-subject",
		},
	})

	resourceKind := "my-resource-kind"
	resourceValue := "my-resource-value"
	action := "my-action"

	expectedErr := errors.New("this is an error returned by compiledManager.IsAllowed()")

	compiledManager := manager.NewMockCompiledPolicy(ctrl)
	compiledManager.EXPECT().IsAllowed("my-subject", resourceKind, resourceValue, action).
		Return(true, expectedErr)

	// When
	result := AuthorizationFunc(compiledManager)(ctx, resourceKind, resourceValue, action)

	// Then
	assert.Equal(t, false, result)
}
