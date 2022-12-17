package attribute

import (
	"errors"
	"regexp"
	"strings"
)

var (
	resourceAttributeRegexp  = regexp.MustCompile(`(resource\.)(.+)`)
	principalAttributeRegexp = regexp.MustCompile(`(principal\.)(.+)`)

	ruleRegexp = regexp.MustCompile(`([resource|principal]?\.?.+)\s?(==|!=)\s?([resource|principal]?\.?.+)`)

	// ErrInvalidRuleFormat is returned when a rule format is invalid.
	ErrInvalidRuleFormat = errors.New("rule is invalid: should have at least one resource.<attribute> or a principal.<attribute>")
)

type RuleOperator string

const (
	// RuleOperatorEqual represents an equal attribute rule.
	// For example: my.owner_id == 123
	RuleOperatorEqual RuleOperator = "=="

	// RuleOperatorEqual represents a NOT equal attribute rule.
	// For example: my.owner_id != 123
	RuleOperatorNotEqual RuleOperator = "!="

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
