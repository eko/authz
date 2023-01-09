package transformer

import (
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/pkg/authz"
)

type principal struct {
	entity *model.Principal
}

func NewPrincipal(entity *model.Principal) *principal {
	return &principal{
		entity: entity,
	}
}

func (t *principal) ToProto() *authz.Principal {
	return &authz.Principal{
		Id:         t.entity.ID,
		Roles:      NewRoles(t.entity.Roles).ToStringSlice(),
		Attributes: NewAttributes(t.entity.Attributes).ToProto(),
	}
}
