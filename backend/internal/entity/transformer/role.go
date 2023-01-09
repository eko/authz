package transformer

import (
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/pkg/authz"
)

type roles struct {
	entities []*model.Role
}

func NewRoles(entities []*model.Role) *roles {
	return &roles{
		entities: entities,
	}
}

func (t *roles) ToProto() []*authz.Role {
	var roles = []*authz.Role{}

	for _, role := range t.entities {
		roles = append(roles, NewRole(role).ToProto())
	}

	return roles
}

func (t *roles) ToStringSlice() []string {
	var roles = []string{}

	for _, role := range t.entities {
		roles = append(roles, NewRole(role).ToString())
	}

	return roles
}

type role struct {
	entity *model.Role
}

func NewRole(entity *model.Role) *role {
	return &role{
		entity: entity,
	}
}

func (t *role) ToProto() *authz.Role {
	return &authz.Role{
		Id:       t.entity.ID,
		Policies: NewPolicies(t.entity.Policies).ToStringSlice(),
	}
}

func (t *role) ToString() string {
	return t.entity.ID
}
