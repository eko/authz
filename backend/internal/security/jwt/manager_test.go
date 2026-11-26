package jwt

import (
	"testing"
	lib_time "time"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/helper/time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestNewManager(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.Auth{}
	clock := time.NewMockClock(ctrl)

	// When
	managerInstance := NewManager(cfg, clock)

	// Then
	assert.IsType(t, new(manager), managerInstance)
}

func TestManager_Generate(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.Auth{
		AccessTokenDuration: 1 * lib_time.Hour,
	}

	date := lib_time.Date(2023, lib_time.January, 1, 0, 0, 0, 0, lib_time.UTC)

	clock := time.NewMockClock(ctrl)
	clock.EXPECT().Now().Return(date).Times(2)

	manager := NewManager(cfg, clock)

	// When
	token, err := manager.Generate("user-123")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, &Token{
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRoeiIsInN1YiI6InVzZXItMTIzIiwiZXhwIjoxNjcyNTM0ODAwLCJpYXQiOjE2NzI1MzEyMDB9.v4KOq5Fw-2PoVPNNMDgunVO4R_RO0NwdLUlyupNXvSk",
		TokenType: "bearer",
		ExpiresIn: 3600,
	}, token)
}

func TestManager_Parse_WhenInvalidToken(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.Auth{
		AccessTokenDuration: 1 * lib_time.Hour,
	}

	date := lib_time.Date(2023, lib_time.January, 1, 0, 0, 0, 0, lib_time.UTC)

	clock := time.NewMockClock(ctrl)
	clock.EXPECT().Now().Return(date).Times(3)

	manager := NewManager(cfg, clock)

	// When
	token, err := manager.Generate("user-123")
	assert.Nil(t, err)

	claims, err := manager.Parse(token.Token)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "user-123", claims.Subject)
}

func TestManager_Parse_WhenSuccess(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.Auth{
		AccessTokenDuration: 1 * lib_time.Hour,
	}

	clock := time.NewMockClock(ctrl)

	manager := NewManager(cfg, clock)

	// When
	claims, err := manager.Parse("this-is-an-invalid-token")

	// Then
	assert.ErrorIs(t, err, jwt.ErrTokenMalformed)
	assert.ErrorContains(t, err, "token contains an invalid number of segments")
	assert.Nil(t, claims)
}

func TestManager_Parse_RejectsForgedTokens(t *testing.T) {
	date := lib_time.Date(2023, lib_time.January, 1, 0, 0, 0, 0, lib_time.UTC)

	claims := jwt.RegisteredClaims{
		Issuer:    configs.ApplicationName,
		Subject:   "authz-user-admin",
		IssuedAt:  jwt.NewNumericDate(date),
		ExpiresAt: jwt.NewNumericDate(date.Add(lib_time.Hour)),
	}

	signedWithOtherKey, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte("4uthz-s3cr3t-valu3-pl3as3-ch4ng3!"))
	assert.NoError(t, err)

	signedWithOtherAlg, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).
		SignedString([]byte("expected-secret-value-of-32-bytes!!"))
	assert.NoError(t, err)

	unsigned, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).
		SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	testCases := []struct {
		name  string
		token string
	}{
		{name: "signed with another key", token: signedWithOtherKey},
		{name: "signed with another algorithm", token: signedWithOtherAlg},
		{name: "unsigned", token: unsigned},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			clock := time.NewMockClock(ctrl)
			clock.EXPECT().Now().Return(date).AnyTimes()

			manager := NewManager(&configs.Auth{
				JWTSignString: "expected-secret-value-of-32-bytes!!",
			}, clock)

			result, err := manager.Parse(testCase.token)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}
