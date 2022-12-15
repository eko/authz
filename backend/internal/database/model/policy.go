package model

import "time"

type Policy struct {
	ID        string      `json:"id" gorm:"primarykey"`
	Resources []*Resource `json:"resources,omitempty" gorm:"many2many:authz_policies_resources;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Actions   []*Action   `json:"actions,omitempty" gorm:"many2many:authz_policies_actions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (Policy) TableName() string {
	return "authz_policies"
}
