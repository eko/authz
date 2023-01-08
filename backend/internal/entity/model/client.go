package model

import (
	"time"
)

type Client struct {
	ID        string    `json:"client_id" gorm:"primarykey"`
	Secret    string    `json:"client_secret" gorm:"type:varchar(512)"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain" gorm:"type:varchar(512)"`
	Data      string    `json:"data,omitempty" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Client) TableName() string {
	return "authz_clients"
}
