package model

import "time"

type Audit struct {
	ID            int64     `json:"id" gorm:"primarykey;autoIncrement"`
	Date          time.Time `json:"date"`
	Principal     string    `json:"principal"`
	ResourceKind  string    `json:"resource_kind"`
	ResourceValue string    `json:"resource_value"`
	Action        string    `json:"action"`
	IsAllowed     bool      `json:"is_allowed"`
	PolicyID      string    `json:"policy_id"`
}

func (Audit) TableName() string {
	return "authz_audit"
}
