package attribute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringToRuleOperator(t *testing.T) {
	testCases := []struct {
		name         string
		rule         string
		expectedRule *Rule
		expectedErr  error
	}{
		{
			name: "valid equal rule",
			rule: "resource.my_attribute==my_value",
			expectedRule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorEqual,
				Value:             "my_value",
			},
			expectedErr: nil,
		},
		{
			name: "valid equal rule (with whitespaces and reversed)",
			rule: "my_value == resource.my_attribute",
			expectedRule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorEqual,
				Value:             "my_value",
			},
			expectedErr: nil,
		},
		{
			name: "valid equal rule (with both resource and principal attributes)",
			rule: "resource.owner_id == principal.owner.id",
			expectedRule: &Rule{
				ResourceAttribute:  "owner_id",
				PrincipalAttribute: "owner.id",
				Operator:           RuleOperatorEqual,
			},
			expectedErr: nil,
		},
		{
			name: "valid not equal rule",
			rule: "resource.my_attribute != principal.another_attribute",
			expectedRule: &Rule{
				ResourceAttribute:  "my_attribute",
				PrincipalAttribute: "another_attribute",
				Operator:           RuleOperatorNotEqual,
			},
			expectedErr: nil,
		},
		{
			name: "valid not equal rule (with principal attribute and value)",
			rule: "true != principal.another_attribute",
			expectedRule: &Rule{
				PrincipalAttribute: "another_attribute",
				Operator:           RuleOperatorNotEqual,
				Value:              "true",
			},
			expectedErr: nil,
		},
		{
			name:         "invalid equal rule",
			rule:         "my_attribute == my_value",
			expectedRule: nil,
			expectedErr:  ErrInvalidRuleFormat,
		},
		{
			name:         "invalid not equal rule",
			rule:         "some.attribute != true",
			expectedRule: nil,
			expectedErr:  ErrInvalidRuleFormat,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rule, err := ConvertStringToRuleOperator(testCase.rule)

			assert.Equal(t, testCase.expectedRule, rule)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestConvertRuleToString(t *testing.T) {
	testCases := []struct {
		name     string
		rule     *Rule
		expected string
	}{
		{
			name: "valid equal rule",
			rule: &Rule{
				ResourceAttribute:  "my_attribute",
				PrincipalAttribute: "another_attribute",
				Operator:           RuleOperatorEqual,
			},
			expected: "resource.my_attribute == principal.another_attribute",
		},
		{
			name: "valid not equal rule",
			rule: &Rule{
				ResourceAttribute:  "my_attribute",
				PrincipalAttribute: "another_attribute",
				Operator:           RuleOperatorNotEqual,
			},
			expected: "resource.my_attribute != principal.another_attribute",
		},
		{
			name: "valid equal rule (resource attribute with a value)",
			rule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorEqual,
				Value:             "something",
			},
			expected: "resource.my_attribute == something",
		},
		{
			name: "valid equal rule (principal attribute with a value)",
			rule: &Rule{
				PrincipalAttribute: "my_attribute",
				Operator:           RuleOperatorEqual,
				Value:              "something",
			},
			expected: "principal.my_attribute == something",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.rule.ToString())
		})
	}
}
