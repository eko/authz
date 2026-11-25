package attribute

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/eko/authz/backend/internal/entity/model"
)

var (
	resourceAttributeRegexp  = regexp.MustCompile(`(resource\.)(.+)`)
	principalAttributeRegexp = regexp.MustCompile(`(principal\.)(.+)`)

	ruleRegexp = regexp.MustCompile(`^\s*([^\s=!<>]+)\s*(==|!=|>=|<=|>|<|contains)\s*(.+?)\s*$`)

	// ErrInvalidRuleFormat is returned when a rule format is invalid.
	ErrInvalidRuleFormat = errors.New("rule is invalid: should have at least one resource.<attribute> or a principal.<attribute>")
)

type RuleOperator string

const (
	// RuleOperatorEqual represents an equal attribute rule.
	// For example: my.owner_id == 123
	RuleOperatorEqual RuleOperator = "=="

	// RuleOperatorGreater represents a greater value attribute rule.
	// For example: my.number > 123
	RuleOperatorGreater RuleOperator = ">"

	// RuleOperatorGreater represents a greater or equal value attribute rule.
	// For example: my.number >= 123
	RuleOperatorGreaterEqual RuleOperator = ">="

	// RuleOperatorLower represents a lower value attribute rule.
	// For example: my.number < 123
	RuleOperatorLower RuleOperator = "<"

	// RuleOperatorLowerEqual represents a lower or equal value attribute rule.
	// For example: my.number <= 123
	RuleOperatorLowerEqual RuleOperator = "<="

	// RuleOperatorEqual represents a NOT equal attribute rule.
	// For example: my.owner_id != 123
	RuleOperatorNotEqual RuleOperator = "!="

	// RuleOperatorContains represents a rule matching when one side, as a
	// comma-separated list, contains the other side's value.
	// For example: principal.applications contains app1
	RuleOperatorContains RuleOperator = "contains"

	principal = "principal"
	resource  = "resource"
)

// Rule represents an attribute rule containing the attribute name and
// the operator to apply to a given value.
type Rule struct {
	ResourceAttribute  string       `json:"resource_attribute"`
	PrincipalAttribute string       `json:"principal_attribute"`
	Operator           RuleOperator `json:"operator"`
	Value              string       `json:"Value"`
}

func (r *Rule) MatchPrincipal(attributes model.Attributes) bool {
	value := attributes.GetAttribute(r.PrincipalAttribute)

	if r.PrincipalAttribute == "" || value == "" {
		return true
	}

	return r.match(value)
}

func (r *Rule) MatchResource(attributes model.Attributes) bool {
	value := attributes.GetAttribute(r.ResourceAttribute)

	if r.ResourceAttribute == "" || value == "" {
		return true
	}

	return r.match(value)
}

func (r *Rule) match(value string) bool {
	return r.Compare(value, r.Value)
}

// Compare applies the rule operator between left and right values.
func (r *Rule) Compare(left string, right string) bool {
	switch r.Operator {
	case RuleOperatorEqual:
		return left == right
	case RuleOperatorNotEqual:
		return left != right
	case RuleOperatorContains:
		return listContains(left, right) || listContains(right, left)
	case RuleOperatorGreater, RuleOperatorGreaterEqual, RuleOperatorLower, RuleOperatorLowerEqual:
		return compareInts(r.Operator, left, right)
	default:
		return false
	}
}

func compareInts(operator RuleOperator, left string, right string) bool {
	leftInt, leftErr := strconv.ParseInt(left, 10, 0)
	rightInt, rightErr := strconv.ParseInt(right, 10, 0)
	if leftErr != nil || rightErr != nil {
		return false
	}

	switch operator {
	case RuleOperatorGreater:
		return leftInt > rightInt
	case RuleOperatorGreaterEqual:
		return leftInt >= rightInt
	case RuleOperatorLower:
		return leftInt < rightInt
	case RuleOperatorLowerEqual:
		return leftInt <= rightInt
	default:
		return false
	}
}

func listContains(list string, value string) bool {
	for _, item := range strings.Split(list, ",") {
		if strings.TrimSpace(item) == value {
			return true
		}
	}

	return false
}

// ToString converts the rule structure to string.
func (r *Rule) ToString() string {
	if (r.ResourceAttribute == "" && r.PrincipalAttribute == "") ||
		(r.ResourceAttribute == "" && r.Value == "") ||
		(r.PrincipalAttribute == "" && r.Value == "") {
		return ""
	}

	if r.ResourceAttribute != "" && r.PrincipalAttribute != "" {
		return resource + "." + r.ResourceAttribute + " " + string(r.Operator) + " " + principal + "." + r.PrincipalAttribute
	} else if r.ResourceAttribute != "" {
		return resource + "." + r.ResourceAttribute + " " + string(r.Operator) + " " + r.Value
	}

	return principal + "." + r.PrincipalAttribute + " " + string(r.Operator) + " " + r.Value
}

// ConvertStringToRuleOperator converts a string to a RuleOperator.
func ConvertStringToRuleOperator(ruleStr string) (*Rule, error) {
	if !ruleRegexp.MatchString(ruleStr) {
		return nil, errors.New("unable to parse attribute rule string")
	}

	ruleMatches := ruleRegexp.FindStringSubmatch(ruleStr)

	value1, operator, value2 := ruleMatches[1], ruleMatches[2], ruleMatches[3]

	var resourceAttribute, principalAttribute, value string

	if resourceAttributeRegexp.MatchString(value1) {
		resourceAttribute = resourceAttributeRegexp.ReplaceAllString(value1, "$2")
	} else if principalAttributeRegexp.MatchString(value1) {
		principalAttribute = principalAttributeRegexp.ReplaceAllString(value1, "$2")
	} else {
		value = value1
	}

	if resourceAttributeRegexp.MatchString(value2) {
		resourceAttribute = resourceAttributeRegexp.ReplaceAllString(value2, "$2")
	} else if principalAttributeRegexp.MatchString(value2) {
		principalAttribute = principalAttributeRegexp.ReplaceAllString(value2, "$2")
	} else {
		value = value2
	}

	if resourceAttribute == "" && principalAttribute == "" {
		return nil, ErrInvalidRuleFormat
	}

	return &Rule{
		ResourceAttribute:  strings.TrimSpace(resourceAttribute),
		PrincipalAttribute: strings.TrimSpace(principalAttribute),
		Operator:           RuleOperator(operator),
		Value:              strings.TrimSpace(value),
	}, nil
}
