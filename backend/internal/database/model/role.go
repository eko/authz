package model

import "time"

type Role struct {
	ID        string    `json:"id" gorm:"primarykey"`
	Policies  []*Policy `json:"policies,omitempty" gorm:"many2many:authz_roles_policies;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Role) TableName() string {
	return "authz_roles"
}
