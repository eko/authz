package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSlugFromString(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "value is empty",
			value:    "",
			expected: true,
		},
		{
			name:     "value is a slug",
			value:    "i-am-a-slug",
			expected: true,
		},
		{
			name:     "value is a slug with special characters",
			value:    "88i-am-a-slug_with.123.special_characters.*",
			expected: true,
		},
		{
			name:     "value is not a slug because of uppercase characters",
			value:    "I-am-a-slug",
			expected: false,
		},
		{
			name:     "value is not a slug",
			value:    "I am not a slug",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := ValidateSlugFromString(testCase.value)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
