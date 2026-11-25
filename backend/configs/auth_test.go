package configs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth_EnsureJWTSignString(t *testing.T) {
	testCases := []struct {
		name          string
		signString    []byte
		expectedErr   error
		expectedValue []byte
	}{
		{
			name:        "legacy default is rejected",
			signString:  []byte(legacyJWTSignString),
			expectedErr: ErrJWTSignStringLegacyDefault,
		},
		{
			name:        "too short secret is rejected",
			signString:  []byte("short-secret"),
			expectedErr: ErrJWTSignStringTooShort,
		},
		{
			name:          "custom strong secret is kept",
			signString:    []byte(strings.Repeat("a", minJWTSignStringLength)),
			expectedValue: []byte(strings.Repeat("a", minJWTSignStringLength)),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := &Auth{JWTSignString: testCase.signString}

			err := cfg.ensureJWTSignString()

			assert.ErrorIs(t, err, testCase.expectedErr)
			if testCase.expectedValue != nil {
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

	assert.Len(t, first.JWTSignString, generatedJWTSignStringLength)
	assert.NotEqual(t, []byte(legacyJWTSignString), first.JWTSignString)
	assert.NotEqual(t, first.JWTSignString, second.JWTSignString)
}
