package transformer

import (
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/pkg/authz"
)

type policies struct {
	entities []*model.Policy
}

func NewPolicies(entities []*model.Policy) *policies {
	return &policies{
		entities: entities,
	}
}

func (t *policies) ToStringSlice() []string {
	var policies = []string{}

	for _, policy := range t.entities {
		policies = append(policies, NewPolicy(policy).ToString())
	}

	return policies
}

func (t *policies) ToProto() []*authz.Policy {
	var policies = []*authz.Policy{}

	for _, policy := range t.entities {
		policies = append(policies, NewPolicy(policy).ToProto())
	}

	return policies
}

type policy struct {
	entity *model.Policy
}

func NewPolicy(entity *model.Policy) *policy {
	return &policy{
		entity: entity,
	}
}

func (t *policy) ToProto() *authz.Policy {
	return &authz.Policy{
		Id:             t.entity.ID,
		Actions:        NewActions(t.entity.Actions).ToStringSlice(),
		Resources:      NewResources(t.entity.Resources).ToStringSlice(),
		AttributeRules: t.entity.AttributeRules.Data,
	}
}

func (t *policy) ToString() string {
	return t.entity.ID
}
