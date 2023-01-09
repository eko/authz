package transformer

import (
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/pkg/authz"
)

type attributes struct {
	entities []*model.Attribute
}

func NewAttributes(entities []*model.Attribute) *attributes {
	return &attributes{
		entities: entities,
	}
}

func (t *attributes) ToProto() []*authz.Attribute {
	var attributes = []*authz.Attribute{}

	for _, attribute := range t.entities {
		attributes = append(attributes, NewAttribute(attribute).ToProto())
	}

	return attributes
}

type attribute struct {
	entity *model.Attribute
}

func NewAttribute(entity *model.Attribute) *attribute {
	return &attribute{
		entity: entity,
	}
}

func (t *attribute) ToProto() *authz.Attribute {
	return &authz.Attribute{
		Key:   t.entity.Key,
		Value: t.entity.Value,
	}
}
