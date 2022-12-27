package model

import "time"

type User struct {
	Username     string    `json:"username" gorm:"primarykey"`
	PasswordHash string    `json:"-" gorm:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "authz_users"
}
