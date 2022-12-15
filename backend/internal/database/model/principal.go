package model

import "time"

type Principal struct {
	ID        int64     `json:"id" gorm:"primarykey;autoIncrement"`
	IsLocked  bool      `json:"is_locked" gorm:"is_locked"`
	Value     string    `json:"value" gorm:"value;uniqueIndex"`
	Roles     []*Role   `json:"roles,omitempty" gorm:"many2many:authz_principals_roles;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Principal) TableName() string {
	return "authz_principals"
}
