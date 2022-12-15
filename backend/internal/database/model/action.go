package model

import "time"

type Action struct {
	ID        string    `json:"id" gorm:"primarykey"`
	IsLocked  bool      `json:"is_locked" gorm:"is_locked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Action) TableName() string {
	return "authz_actions"
}
