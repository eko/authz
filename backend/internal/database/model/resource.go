package model

import (
	"time"
)

type Resource struct {
	ID         string       `json:"id" gorm:"primarykey"`
	Kind       string       `json:"kind" gorm:"kind"`
	Value      string       `json:"value" gorm:"value"`
	Attributes []*Attribute `json:"attributes,omitempty" gorm:"many2many:authz_resources_attributes;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsLocked   bool         `json:"is_locked" gorm:"is_locked"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

func (Resource) TableName() string {
	return "authz_resources"
}

func (r *Resource) GetAttribute(key string) string {
	for _, attribute := range r.Attributes {
		if attribute.Key == key {
			return attribute.Value
		}
	}

	return ""
}
