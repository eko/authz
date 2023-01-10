package interceptor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eko/authz/backend/internal/security/jwt"
	lib_jwt "github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthenticateFunc_WhenSuccess(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	md := metadata.New(map[string]string{
		"authorization": "bearer token-123",
	})

	ctx = metadata.NewIncomingContext(ctx, md)

	expectedClaims := &jwt.Claims{
		RegisteredClaims: lib_jwt.RegisteredClaims{
			Issuer:    "authz",
			Subject:   "someone",
			ExpiresAt: lib_jwt.NewNumericDate(time.Date(2023, time.January, 2, 8, 0, 0, 0, time.UTC)),
		},
	}

	tokenManager := jwt.NewMockManager(ctrl)
	tokenManager.EXPECT().Parse("token-123").Return(expectedClaims, nil)

	// When
	newCtx, err := AuthenticateFunc(tokenManager)(ctx)

	// Then
	assert.Nil(t, err)

	claims := newCtx.Value(ClaimsKey)
	assert.Equal(t, expectedClaims, claims)
}

func TestAuthenticateFunc_WhenParseError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	md := metadata.New(map[string]string{
		"authorization": "bearer token-123",
	})

	ctx = metadata.NewIncomingContext(ctx, md)

	tokenManager := jwt.NewMockManager(ctrl)

	expectedErr := errors.New("this is an expected error from test")
	tokenManager.EXPECT().Parse("token-123").Return(nil, expectedErr)

	// When
	newCtx, err := AuthenticateFunc(tokenManager)(ctx)

	// Then
	assert.Nil(t, newCtx)
	assert.Equal(t, err, status.Errorf(codes.Unauthenticated, "unable to parse token: %v", expectedErr))
}
