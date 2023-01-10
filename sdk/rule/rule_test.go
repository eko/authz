package rule

import "testing"

func TestPrincipalAttributeEqualValue(t *testing.T) {
	testCases := []struct {
		name      string
		attribute string
		value     string
		expected  string
	}{
		{
			name:      "empty values",
			attribute: "",
			value:     "",
			expected:  "",
		},
		{
			name:      "attribute value empty",
			attribute: "",
			value:     "something",
			expected:  "",
		},
		{
			name:      "value empty",
			attribute: "something",
			value:     "",
			expected:  "",
		},
		{
			name:      "both attribute and value filled",
			attribute: "something",
			value:     "my_value",
			expected:  "principal.something == my_value",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value := PrincipalAttributeEqualValue(testCase.attribute, testCase.value)

			if value != testCase.expected {
				t.Fatalf("unexpected value received: %s, expected: %s", value, testCase.expected)
			}
		})
	}
}

func TestPrincipalAttributeNotEqualValue(t *testing.T) {
	testCases := []struct {
		name      string
		attribute string
		value     string
		expected  string
	}{
		{
			name:      "empty values",
			attribute: "",
			value:     "",
			expected:  "",
		},
		{
			name:      "attribute value empty",
			attribute: "",
			value:     "something",
			expected:  "",
		},
		{
			name:      "value empty",
			attribute: "something",
			value:     "",
			expected:  "",
		},
		{
			name:      "both attribute and value filled",
			attribute: "something",
			value:     "my_value",
			expected:  "principal.something != my_value",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value := PrincipalAttributeNotEqualValue(testCase.attribute, testCase.value)

			if value != testCase.expected {
				t.Fatalf("unexpected value received: %s, expected: %s", value, testCase.expected)
			}
		})
	}
}

func TestResourceAttributeEqualValue(t *testing.T) {
	testCases := []struct {
		name      string
		attribute string
		value     string
		expected  string
	}{
		{
			name:      "empty values",
			attribute: "",
			value:     "",
			expected:  "",
		},
		{
			name:      "attribute value empty",
			attribute: "",
			value:     "something",
			expected:  "",
		},
		{
			name:      "value empty",
			attribute: "something",
			value:     "",
			expected:  "",
		},
		{
			name:      "both attribute and value filled",
			attribute: "something",
			value:     "my_value",
			expected:  "resource.something == my_value",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value := ResourceAttributeEqualValue(testCase.attribute, testCase.value)

			if value != testCase.expected {
				t.Fatalf("unexpected value received: %s, expected: %s", value, testCase.expected)
			}
		})
	}
}

func TestResourceAttributeNotEqualValue(t *testing.T) {
	testCases := []struct {
		name      string
		attribute string
		value     string
		expected  string
	}{
		{
			name:      "empty values",
			attribute: "",
			value:     "",
			expected:  "",
		},
		{
			name:      "attribute value empty",
			attribute: "",
			value:     "something",
			expected:  "",
		},
		{
			name:      "value empty",
			attribute: "something",
			value:     "",
			expected:  "",
		},
		{
			name:      "both attribute and value filled",
			attribute: "something",
			value:     "my_value",
			expected:  "resource.something != my_value",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value := ResourceAttributeNotEqualValue(testCase.attribute, testCase.value)

			if value != testCase.expected {
				t.Fatalf("unexpected value received: %s, expected: %s", value, testCase.expected)
			}
		})
	}
}

func TestAttributeEqualValue(t *testing.T) {
	testCases := []struct {
		name       string
		attributes PrincipalResourceAttribute
		expected   string
	}{
		{
			name:       "empty values",
			attributes: PrincipalResourceAttribute{},
			expected:   "",
		},
		{
			name: "principal attribute value empty",
			attributes: PrincipalResourceAttribute{
				ResourceAttribute: "something",
			},
			expected: "",
		},
		{
			name: "resource attribute value empty",
			attributes: PrincipalResourceAttribute{
				PrincipalAttribute: "something",
			},
			expected: "",
		},
		{
			name: "both attribute and value filled",
			attributes: PrincipalResourceAttribute{
				PrincipalAttribute: "email",
				ResourceAttribute:  "owner_email",
			},
			expected: "principal.email == resource.owner_email",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value := AttributeEqual(testCase.attributes)

			if value != testCase.expected {
				t.Fatalf("unexpected value received: %s, expected: %s", value, testCase.expected)
			}
		})
	}
}

func TestAttributeNotEqualValue(t *testing.T) {
	testCases := []struct {
		name       string
		attributes PrincipalResourceAttribute
		expected   string
	}{
		{
			name:       "empty values",
			attributes: PrincipalResourceAttribute{},
			expected:   "",
		},
		{
			name: "principal attribute value empty",
			attributes: PrincipalResourceAttribute{
				ResourceAttribute: "something",
			},
			expected: "",
		},
		{
			name: "resource attribute value empty",
			attributes: PrincipalResourceAttribute{
				PrincipalAttribute: "something",
			},
			expected: "",
		},
		{
			name: "both attribute and value filled",
			attributes: PrincipalResourceAttribute{
				PrincipalAttribute: "email",
				ResourceAttribute:  "owner_email",
			},
			expected: "principal.email != resource.owner_email",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value := AttributeNotEqual(testCase.attributes)

			if value != testCase.expected {
				t.Fatalf("unexpected value received: %s, expected: %s", value, testCase.expected)
			}
		})
	}
}
