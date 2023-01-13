package transformer

import (
	"testing"

	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/pkg/authz"
	"github.com/stretchr/testify/assert"
)

func TestNewAttribute_ToProto(t *testing.T) {
	// Given
	attribute := &model.Attribute{
		ID:    1,
		Key:   "key1",
		Value: "value1",
	}

	// When
	result := NewAttribute(attribute).ToProto()

	// Then
	assert.IsType(t, new(authz.Attribute), result)

	assert.Equal(t, "key1", result.Key)
	assert.Equal(t, "value1", result.Value)
}

func TestNewAttributes_ToProto(t *testing.T) {
	// Given
	attribute1 := &model.Attribute{
		ID:    1,
		Key:   "key1",
		Value: "value1",
	}

	attribute2 := &model.Attribute{
		ID:    2,
		Key:   "key2",
		Value: "value2",
	}

	// When
	result := NewAttributes([]*model.Attribute{
		attribute1,
		attribute2,
	}).ToProto()

	// Then
	assert.Equal(t, "key1", result[0].Key)
	assert.Equal(t, "value1", result[0].Value)

	assert.Equal(t, "key2", result[1].Key)
	assert.Equal(t, "value2", result[1].Value)
}
