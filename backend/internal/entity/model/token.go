package model

import "time"

type Token struct {
	ID        uint   `gorm:"primarykey"`
	Code      string `gorm:"type:varchar(512)"`
	Access    string `gorm:"type:varchar(512)"`
	Refresh   string `gorm:"type:varchar(512)"`
	Data      string `gorm:"type:text"`
	ExpiredAt int64
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Token) TableName() string {
	return "authz_oauth_tokens"
}
