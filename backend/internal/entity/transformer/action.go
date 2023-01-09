package transformer

import (
	"github.com/eko/authz/backend/internal/entity/model"
)

type actions struct {
	entities []*model.Action
}

func NewActions(entities []*model.Action) *actions {
	return &actions{
		entities: entities,
	}
}

func (t *actions) ToStringSlice() []string {
	var actions = []string{}

	for _, action := range t.entities {
		actions = append(actions, NewAction(action).ToString())
	}

	return actions
}

type action struct {
	entity *model.Action
}

func NewAction(entity *model.Action) *action {
	return &action{
		entity: entity,
	}
}

func (t *action) ToString() string {
	return t.entity.ID
}
