package model

import "time"

type Action struct {
	ID        int64     `json:"id" gorm:"primarykey;autoIncrement"`
	IsLocked  bool      `json:"is_locked" gorm:"is_locked"`
	Name      string    `json:"name" gorm:"name;uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Action) TableName() string {
	return "authz_actions"
}
