package model

import (
	"time"
)

type Principal struct {
	ID         string       `json:"id" gorm:"primarykey"`
	Roles      []*Role      `json:"roles,omitempty" gorm:"many2many:authz_principals_roles;"`
	Attributes []*Attribute `json:"attributes,omitempty" gorm:"many2many:authz_principals_attributes;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsLocked   bool         `json:"is_locked" gorm:"is_locked"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

func (Principal) TableName() string {
	return "authz_principals"
}
