package model

import "time"

type CompiledPolicy struct {
	PolicyID      string    `gorm:"index"`
	ResourceKind  string    `gorm:"index"`
	ResourceValue string    `gorm:"index"`
	ActionID      string    `gorm:"index"`
	Version       int64     `gorm:"index"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (CompiledPolicy) TableName() string {
	return "authz_compiled_policies"
}
