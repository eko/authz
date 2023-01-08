package manager

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// ResourceSeparator is the resource name separator between kind and value.
	ResourceSeparator = "."

	// WildcardValue is the wildcard value used to identify resources.
	WildcardValue = "*"
)

// ResourceSplit splits a resource name string to:
// * resource kind,
// * resource value.
func ResourceSplit(resource string) (string, string) {
	parts := strings.Split(resource, ResourceSeparator)

	kind := strings.Join(parts[:len(parts)-1], ResourceSeparator)
	value := strings.Join(parts[len(parts)-1:], ResourceSeparator)

	return kind, value
}

// CastAnyToString casts a value of any type to a string.
func CastAnyToString(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("unsupported attribute type: %T", value)
	}
}
