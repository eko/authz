package transformer

import (
	"testing"

	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/pkg/authz"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestNewPolicy_ToProto(t *testing.T) {
	// Given
	attribute := &model.Policy{
		ID: "policy-1",
		Resources: []*model.Resource{
			{
				ID:    "resource-1",
				Kind:  "kind-1",
				Value: "value-1",
				Attributes: []*model.Attribute{
					{ID: 1, Key: "key1", Value: "value1"},
				},
				IsLocked: false,
			},
		},
		Actions: []*model.Action{
			{
				ID: "action-1",
			},
		},
		AttributeRules: datatypes.NewJSONType(
			[]string{"rule1", "rule2"},
		),
	}

	// When
	result := NewPolicy(attribute).ToProto()

	// Then
	assert.IsType(t, new(authz.Policy), result)

	assert.Equal(t, "policy-1", result.Id)

	assert.Equal(t, "resource-1", result.Resources[0])
	assert.Equal(t, "action-1", result.Actions[0])

	assert.Equal(t, "rule1", result.AttributeRules[0])
	assert.Equal(t, "rule2", result.AttributeRules[1])
}

func TestNewPolicys_ToProto(t *testing.T) {
	// Given
	attribute1 := &model.Policy{
		ID: "policy-1",
		Resources: []*model.Resource{
			{
				ID:    "resource-1",
				Kind:  "kind-1",
				Value: "value-1",
				Attributes: []*model.Attribute{
					{ID: 1, Key: "key1", Value: "value1"},
				},
				IsLocked: false,
			},
		},
		Actions: []*model.Action{
			{
				ID: "action-1",
			},
		},
		AttributeRules: datatypes.NewJSONType(
			[]string{"rule1", "rule2"},
		),
	}

	attribute2 := &model.Policy{
		ID: "policy-2",
		Resources: []*model.Resource{
			{
				ID:    "resource-2",
				Kind:  "kind2",
				Value: "value2",
				Attributes: []*model.Attribute{
					{ID: 1, Key: "key2", Value: "value2"},
				},
				IsLocked: false,
			},
		},
		Actions: []*model.Action{
			{
				ID: "action-2",
			},
		},
		AttributeRules: datatypes.NewJSONType(
			[]string{"rule1", "rule2"},
		),
	}

	// When
	result := NewPolicies([]*model.Policy{
		attribute1,
		attribute2,
	}).ToProto()

	// Then
	assert.IsType(t, []*authz.Policy{}, result)

	assert.Equal(t, "policy-1", result[0].Id)
	assert.Equal(t, "resource-1", result[0].Resources[0])
	assert.Equal(t, "action-1", result[0].Actions[0])
	assert.Equal(t, "rule1", result[0].AttributeRules[0])
	assert.Equal(t, "rule2", result[0].AttributeRules[1])

	assert.Equal(t, "policy-2", result[1].Id)
	assert.Equal(t, "resource-2", result[1].Resources[0])
	assert.Equal(t, "action-2", result[1].Actions[0])
	assert.Equal(t, "rule1", result[1].AttributeRules[0])
	assert.Equal(t, "rule2", result[1].AttributeRules[1])
}
