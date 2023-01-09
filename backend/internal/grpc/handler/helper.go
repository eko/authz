package handler

import "github.com/eko/authz/backend/pkg/authz"

func attributesMap(attributes []*authz.Attribute) map[string]any {
	var result = map[string]any{}

	for _, attribute := range attributes {
		result[attribute.GetKey()] = attribute.GetValue()
	}

	return result
}
