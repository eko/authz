package configs

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth_EnsureJWTSignString(t *testing.T) {
	testCases := []struct {
		name          string
		signString    string
		expectedErr   error
		expectedValue string
	}{
		{
			name:        "legacy default is rejected",
			signString:  legacyJWTSignString,
			expectedErr: ErrJWTSignStringLegacyDefault,
		},
		{
			name:        "too short secret is rejected",
			signString:  "short-secret",
			expectedErr: ErrJWTSignStringTooShort,
		},
		{
			name:          "custom strong secret is kept",
			signString:    strings.Repeat("a", minJWTSignStringLength),
			expectedValue: strings.Repeat("a", minJWTSignStringLength),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := &Auth{JWTSignString: testCase.signString}

			err := cfg.ensureJWTSignString()

			assert.ErrorIs(t, err, testCase.expectedErr)
			if testCase.expectedValue != "" {
				assert.Equal(t, testCase.expectedValue, cfg.JWTSignString)
			}
		})
	}
}

func TestAuth_EnsureJWTSignString_GeneratesRandomSecretWhenEmpty(t *testing.T) {
	first := newAuth()
	second := newAuth()

	assert.NoError(t, first.ensureJWTSignString())
	assert.NoError(t, second.ensureJWTSignString())

	assert.Len(t, first.JWTSignString, 64)
	assert.NotEqual(t, legacyJWTSignString, first.JWTSignString)
	assert.NotEqual(t, first.JWTSignString, second.JWTSignString)
}

func TestLoad_ReadsJWTSignStringFromEnvironment(t *testing.T) {
	secret := "a-strong-signing-secret-from-environment"
	t.Setenv("AUTH_JWT_SIGN_STRING", secret)

	cfg := Load(context.Background())

	assert.Equal(t, secret, cfg.Auth.JWTSignString)
}
