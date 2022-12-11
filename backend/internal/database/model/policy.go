package model

import "time"

type Policy struct {
	ID        int64       `json:"id" gorm:"primarykey;autoIncrement"`
	Name      string      `json:"name" gorm:"name;uniqueIndex"`
	Resources []*Resource `json:"resources,omitempty" gorm:"many2many:authz_policies_resources;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Actions   []*Action   `json:"actions,omitempty" gorm:"many2many:authz_policies_actions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (Policy) TableName() string {
	return "authz_policies"
}
