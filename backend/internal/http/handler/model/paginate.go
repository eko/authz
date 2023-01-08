package model

import "github.com/eko/authz/backend/internal/entity/model"

type Paginated[T model.Models] struct {
	Data  []*T  `json:"data"`
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Size  int64 `json:"size"`
}

func NewPaginated[T model.Models](data []*T, total, page, size int64) *Paginated[T] {
	return &Paginated[T]{
		Data:  data,
		Total: total,
		Page:  page,
		Size:  size,
	}
}
