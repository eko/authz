package rule

import "fmt"

const (
	// PrincipalKey is the key used in Authz to identify a principal.
	PrincipalKey = "principal"
	// ResourceKey is the key used in Authz to identify a resource.
	ResourceKey = "resource"
)

// PrincipalResourceAttribute is used by AttributeEqual and AttributeNotEqual functions
// to specify both principal attribute key and resource attribute key.
type PrincipalResourceAttribute struct {
	PrincipalAttribute string
	ResourceAttribute  string
}

// PrincipalAttributeEqualValue is used to create an attribute rule when a
// principal attribute value is equal to the given value.
func PrincipalAttributeEqualValue(attribute string, value string) string {
	if attribute == "" || value == "" {
		return ""
	}

	return fmt.Sprintf("%s.%s %s %s", PrincipalKey, attribute, OperatorEqual, value)
}

// PrincipalAttributeNotEqualValue is used to create an attribute rule when a
// principal attribute value is not equal to the given value.
func PrincipalAttributeNotEqualValue(attribute string, value string) string {
	if attribute == "" || value == "" {
		return ""
	}

	return fmt.Sprintf("%s.%s %s %s", PrincipalKey, attribute, OperatorNotEqual, value)
}

// ResourceAttributeEqualValue is used to create an attribute rule when a
// resource attribute value is equal to the given value.
func ResourceAttributeEqualValue(attribute string, value string) string {
	if attribute == "" || value == "" {
		return ""
	}

	return fmt.Sprintf("%s.%s %s %s", ResourceKey, attribute, OperatorEqual, value)
}

// ResourceAttributeNotEqualValue is used to create an attribute rule when a
// resource attribute value is not equal to the given value.
func ResourceAttributeNotEqualValue(attribute string, value string) string {
	if attribute == "" || value == "" {
		return ""
	}

	return fmt.Sprintf("%s.%s %s %s", ResourceKey, attribute, OperatorNotEqual, value)
}

// AttributeEqual is used to create an attribute rule when a
// principal attribute value is equal to the value of given resource attribute value.
func AttributeEqual(value PrincipalResourceAttribute) string {
	if value.PrincipalAttribute == "" || value.ResourceAttribute == "" {
		return ""
	}

	return fmt.Sprintf("%s.%s %s %s.%s", PrincipalKey, value.PrincipalAttribute, OperatorEqual, ResourceKey, value.ResourceAttribute)
}

// AttributeNotEqual is used to create an attribute rule when a
// principal attribute value is not equal to the value of given resource attribute value.
func AttributeNotEqual(value PrincipalResourceAttribute) string {
	if value.PrincipalAttribute == "" || value.ResourceAttribute == "" {
		return ""
	}

	return fmt.Sprintf("%s.%s %s %s.%s", PrincipalKey, value.PrincipalAttribute, OperatorNotEqual, ResourceKey, value.ResourceAttribute)
}
