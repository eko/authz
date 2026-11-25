package attribute

import (
	"testing"

	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/stretchr/testify/assert"
)

func TestRule_MatchPrincipal(t *testing.T) {
	// Given
	testCases := []struct {
		name       string
		attributes model.Attributes
		rule       *Rule
		expected   bool
	}{
		{
			name:       "Empty attributes and rules",
			attributes: model.Attributes{},
			rule:       &Rule{},
			expected:   true,
		},
		{
			name: "Attributes matching > rule",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				PrincipalAttribute: "my_attribute",
				Operator:           RuleOperatorGreater,
				Value:              "40",
			},
			expected: true,
		},
		{
			name: "Attributes not matching >= rule when equal",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				PrincipalAttribute: "my_attribute",
				Operator:           RuleOperatorGreaterEqual,
				Value:              "50",
			},
			expected: true,
		},
		{
			name: "Attributes not matching >= rule when greater",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "51"},
			},
			rule: &Rule{
				PrincipalAttribute: "my_attribute",
				Operator:           RuleOperatorGreaterEqual,
				Value:              "50",
			},
			expected: true,
		},
		{
			name: "Attributes not matching > rule",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				PrincipalAttribute: "my_attribute",
				Operator:           RuleOperatorGreater,
				Value:              "60",
			},
			expected: false,
		},
		{
			name: "Attributes matching < rule",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				PrincipalAttribute: "my_attribute",
				Operator:           RuleOperatorLower,
				Value:              "60",
			},
			expected: true,
		},
		{
			name: "Attributes matching <= rule when equal",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				PrincipalAttribute: "my_attribute",
				Operator:           RuleOperatorLowerEqual,
				Value:              "50",
			},
			expected: true,
		},
		{
			name: "Attributes matching <= rule when lower",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "49"},
			},
			rule: &Rule{
				PrincipalAttribute: "my_attribute",
				Operator:           RuleOperatorLowerEqual,
				Value:              "50",
			},
			expected: true,
		},
		{
			name: "Attributes not matching < rule",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				PrincipalAttribute: "my_attribute",
				Operator:           RuleOperatorLower,
				Value:              "40",
			},
			expected: false,
		},
	}

	// When - Then
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.rule.MatchPrincipal(testCase.attributes))
		})
	}
}

func TestRule_MatchResource(t *testing.T) {
	// Given
	testCases := []struct {
		name       string
		attributes model.Attributes
		rule       *Rule
		expected   bool
	}{
		{
			name:       "Empty attributes and rules",
			attributes: model.Attributes{},
			rule:       &Rule{},
			expected:   true,
		},
		{
			name: "Attributes matching > rule",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorGreater,
				Value:             "40",
			},
			expected: true,
		},
		{
			name: "Attributes not matching >= rule when equal",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorGreaterEqual,
				Value:             "50",
			},
			expected: true,
		},
		{
			name: "Attributes not matching >= rule when greater",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "51"},
			},
			rule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorGreaterEqual,
				Value:             "50",
			},
			expected: true,
		},
		{
			name: "Attributes not matching > rule",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorGreater,
				Value:             "60",
			},
			expected: false,
		},
		{
			name: "Attributes matching < rule",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorLower,
				Value:             "60",
			},
			expected: true,
		},
		{
			name: "Attributes matching <= rule when equal",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorLowerEqual,
				Value:             "50",
			},
			expected: true,
		},
		{
			name: "Attributes matching <= rule when lower",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "49"},
			},
			rule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorLowerEqual,
				Value:             "50",
			},
			expected: true,
		},
		{
			name: "Attributes not matching < rule",
			attributes: model.Attributes{
				{Key: "my_attribute", Value: "50"},
			},
			rule: &Rule{
				ResourceAttribute: "my_attribute",
				Operator:          RuleOperatorLower,
				Value:             "40",
			},
			expected: false,
		},
	}

	// When - Then
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.rule.MatchResource(testCase.attributes))
		})
	}
}

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
			name: "valid greater or equal rule",
			rule: "principal.age >= 18",
			expectedRule: &Rule{
				PrincipalAttribute: "age",
				Operator:           RuleOperatorGreaterEqual,
				Value:              "18",
			},
			expectedErr: nil,
		},
		{
			name: "valid lower or equal rule without whitespaces",
			rule: "resource.level<=principal.level",
			expectedRule: &Rule{
				ResourceAttribute:  "level",
				PrincipalAttribute: "level",
				Operator:           RuleOperatorLowerEqual,
			},
			expectedErr: nil,
		},
		{
			name: "valid contains rule",
			rule: "principal.applications contains app1",
			expectedRule: &Rule{
				PrincipalAttribute: "applications",
				Operator:           RuleOperatorContains,
				Value:              "app1",
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

func TestRule_Compare(t *testing.T) {
	testCases := []struct {
		name     string
		operator RuleOperator
		left     string
		right    string
		expected bool
	}{
		{name: "equal", operator: RuleOperatorEqual, left: "a", right: "a", expected: true},
		{name: "not equal", operator: RuleOperatorNotEqual, left: "a", right: "a", expected: false},
		{name: "greater or equal", operator: RuleOperatorGreaterEqual, left: "18", right: "18", expected: true},
		{name: "lower or equal", operator: RuleOperatorLowerEqual, left: "19", right: "18", expected: false},
		{name: "greater with non numeric", operator: RuleOperatorGreater, left: "abc", right: "1", expected: false},
		{name: "contains in left list", operator: RuleOperatorContains, left: "app1, app2", right: "app2", expected: true},
		{name: "contains in right list", operator: RuleOperatorContains, left: "app1", right: "app2,app1", expected: true},
		{name: "contains does not match prefix", operator: RuleOperatorContains, left: "app10,app2", right: "app1", expected: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rule := &Rule{Operator: testCase.operator}

			assert.Equal(t, testCase.expected, rule.Compare(testCase.left, testCase.right))
		})
	}
}
