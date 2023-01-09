package transformer

import (
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/pkg/authz"
)

type resources struct {
	entities []*model.Resource
}

func NewResources(entities []*model.Resource) *resources {
	return &resources{
		entities: entities,
	}
}

func (t *resources) ToStringSlice() []string {
	var resources = []string{}

	for _, resource := range t.entities {
		resources = append(resources, NewResource(resource).ToString())
	}

	return resources
}

type resource struct {
	entity *model.Resource
}

func NewResource(entity *model.Resource) *resource {
	return &resource{
		entity: entity,
	}
}

func (t *resource) ToProto() *authz.Resource {
	return &authz.Resource{
		Id:         t.entity.ID,
		Kind:       t.entity.Kind,
		Value:      t.entity.Value,
		Attributes: NewAttributes(t.entity.Attributes).ToProto(),
	}
}

func (t *resource) ToString() string {
	return t.entity.ID
}
