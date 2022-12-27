package model

import (
	"time"

	"github.com/lib/pq"
)

type Policy struct {
	ID             string         `json:"id" gorm:"primarykey"`
	Resources      []*Resource    `json:"resources,omitempty" gorm:"many2many:authz_policies_resources;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Actions        []*Action      `json:"actions,omitempty" gorm:"many2many:authz_policies_actions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AttributeRules pq.StringArray `json:"attribute_rules,omitempty" gorm:"type:text[]" swaggertype:"object"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`

	Roles []*Role `json:"-" gorm:"many2many:authz_roles_policies"`
}

func (Policy) TableName() string {
	return "authz_policies"
}
