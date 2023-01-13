package transformer

import (
	"testing"

	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/stretchr/testify/assert"
)

func TestNewPrincipal_ToProto(t *testing.T) {
	// Given
	principal := &model.Principal{
		ID: "principal-1",
		Roles: []*model.Role{
			{
				ID: "role-1",
				Policies: []*model.Policy{
					{
						ID: "policy-1",
						Resources: []*model.Resource{
							{
								ID:    "resource-1",
								Kind:  "kind-1",
								Value: "value-1",
								Attributes: []*model.Attribute{
									{ID: 1, Key: "key1", Value: "value1"},
								},
							},
						},
					},
				},
			},
		},
		Attributes: []*model.Attribute{
			{
				ID:    1,
				Key:   "key1",
				Value: "value1",
			},
		},
	}

	// When
	result := NewPrincipal(principal).ToProto()

	// Then
	assert.Equal(t, "principal-1", result.Id)
	assert.Equal(t, NewRoles(principal.Roles).ToStringSlice(), result.Roles)
	assert.Equal(t, NewAttributes(principal.Attributes).ToProto(), result.Attributes)
}
