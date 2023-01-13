package transformer

import (
	"testing"

	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/stretchr/testify/assert"
)

func TestNewAction_ToString(t *testing.T) {
	// Given
	action := &model.Action{
		ID: "test",
	}

	// When
	result := NewAction(action).ToString()

	// Then
	assert.Equal(t, "test", result)
}

func TestNewActions_ToStringSlice(t *testing.T) {
	// Given
	action1 := &model.Action{
		ID: "action-1",
	}

	action2 := &model.Action{
		ID: "action-2",
	}

	// When
	result := NewActions([]*model.Action{
		action1,
		action2,
	}).ToStringSlice()

	// Then
	assert.Equal(t, []string{"action-1", "action-2"}, result)
}
