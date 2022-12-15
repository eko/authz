package model

import "time"

type Resource struct {
	ID        string    `json:"id" gorm:"primarykey"`
	Kind      string    `json:"kind" gorm:"kind"`
	Value     string    `json:"value" gorm:"value"`
	IsLocked  bool      `json:"is_locked" gorm:"is_locked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Resource) TableName() string {
	return "authz_resources"
}
