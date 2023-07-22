package repository

import (
	"fmt"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/entity/model"
	"gorm.io/gorm"
)

type FieldValue struct {
	Operator string
	Value    any
	Raw      any
}

// QueryOption specifies how options should be formatted.
//
// An option is a function that takes this private struct
// and override values on it.
type QueryOption func(*queryOptions)

type queryOptions struct {
	joins          []string
	preloads       []string
	page           int64
	size           int64
	filter         map[string]FieldValue
	sort           string
	skipPagination bool
}

func WithJoin(joins ...string) QueryOption {
	return func(o *queryOptions) {
		o.joins = joins
	}
}

// WithPreloads allows to specify relationships you want to preload with the
// currently requested resource.
func WithPreloads(preloads ...string) QueryOption {
	return func(o *queryOptions) {
		o.preloads = preloads
	}
}

// WithPage allows to specify the page number you want to retrieve.
func WithPage(page int64) QueryOption {
	return func(o *queryOptions) {
		o.page = page
	}
}

// WithSize allows to specify the size (number of elements) you want to retrieve.
func WithSize(size int64) QueryOption {
	return func(o *queryOptions) {
		o.size = size
	}
}

// WithFilter allows to specify some filters to apply on the query.
func WithFilter(filter map[string]FieldValue) QueryOption {
	return func(o *queryOptions) {
		o.filter = filter
	}
}

// WithSort allows to specify the sort order you want to apply on the query.
func WithSort(sort string) QueryOption {
	return func(o *queryOptions) {
		o.sort = sort
	}
}

// WithSkipPagination allows to skip pagination and retrieve all elements.
func WithSkipPagination() QueryOption {
	return func(o *queryOptions) {
		o.skipPagination = true
	}
}

type Base[T model.Models] interface {
	Create(object ...*T) error
	DB() *gorm.DB
	CountByFields(fieldValues map[string]FieldValue, options ...QueryOption) (int64, error)
	Delete(object *T) error
	DeleteByFields(fieldValues map[string]FieldValue) error
	Find(options ...QueryOption) ([]*T, int64, error)
	Get(identifier string, options ...QueryOption) (*T, error)
	GetByFields(fieldValues map[string]FieldValue, options ...QueryOption) (*T, error)
	Update(object *T) error
	UpdateAssociation(object *T, associationName string, data any) error
	WithTransaction(transaction database.Transaction) *base[T]
}

// base struct that allows contacting the database using Gorm.
type base[T model.Models] struct {
	db *gorm.DB
}

// New initializes a new repository.
func New[T model.Models](db *gorm.DB) Base[T] {
	return &base[T]{
		db: db,
	}
}

// WithTransaction returns a new repository instance using a transaction database.
func (r *base[T]) WithTransaction(transaction database.Transaction) *base[T] {
	return &base[T]{
		db: transaction.DB(),
	}
}

// Create allows to create a new entry in a database table.
func (r *base[T]) Create(object ...*T) error {
	return r.db.Create(object).Error
}

// DB allows returning database session object.
func (r *base[T]) DB() *gorm.DB {
	return r.db
}

// Delete allows to delete the specified entry from the database.
func (r *base[T]) Delete(object *T) error {
	return r.db.Delete(object).Error
}

// DeleteByFields allows to delete values of the current type from the database
// filtered by the given field name and value.
func (r *base[T]) DeleteByFields(fieldValues map[string]FieldValue) error {
	result := new(T)

	db := r.db

	for field, value := range fieldValues {
		if value.Raw != nil {
			db = db.Where(value.Raw)
		} else {
			db = db.Where(fmt.Sprintf("%s %s ?", field, value.Operator), value.Value)
		}
	}

	return db.Delete(result).Error
}

// Update allows to update the specified entry into the database.
func (r *base[T]) Update(object *T) error {
	return r.db.Save(object).Error
}

// UpdateAssociation allows to update the specified association entry into the database.
func (r *base[T]) UpdateAssociation(object *T, associationName string, data any) error {
	return r.db.Model(object).Association(associationName).Replace(data)
}

// Get allows to retrieve a value of the current type from the specified primary key.
func (r *base[T]) Get(identifier string, options ...QueryOption) (*T, error) {
	result := new(T)

	db := r.applyOptions(options)

	if err := db.First(result, "id = ?", identifier).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// Count allows to count values of the current type from the specified field and value.
func (r *base[T]) CountByFields(fieldValues map[string]FieldValue, options ...QueryOption) (int64, error) {
	var total int64
	result := new(T)

	db := r.applyOptions(options)

	for field, value := range fieldValues {
		if value.Raw != nil {
			db = db.Where(value.Raw)
		} else {
			db = db.Where(fmt.Sprintf("%s %s ?", field, value.Operator), value.Value)
		}
	}

	if err := db.Model(&result).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// GetByFields allows to retrieve a value of the current type from the database
// filtered by the given field names and values.
func (r *base[T]) GetByFields(fieldValues map[string]FieldValue, options ...QueryOption) (*T, error) {
	result := new(T)

	db := r.applyOptions(options)

	for field, value := range fieldValues {
		if value.Raw != nil {
			db = db.Where(value.Raw)
		} else {
			db = db.Where(fmt.Sprintf("%s %s ?", field, value.Operator), value.Value)
		}
	}

	if err := db.First(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// Find allows to retrieve a list of values of the current type.
func (r *base[T]) Find(options ...QueryOption) ([]*T, int64, error) {
	var result = make([]*T, 0)
	var total int64

	db := r.applyOptions(options)

	if err := db.Find(&result).Error; err != nil {
		return []*T{}, 0, err
	}

	options = append(options, WithSkipPagination())
	db = r.applyOptions(options)

	if err := db.Model(&result).Count(&total).Error; err != nil {
		return []*T{}, 0, err
	}

	return result, total, nil
}

func (r *base[T]) applyOptions(options []QueryOption) *gorm.DB {
	db := r.db

	opts := &queryOptions{}

	for _, option := range options {
		option(opts)
	}

	if joins := opts.joins; len(joins) > 0 {
		for _, join := range joins {
			db = db.Joins(join)
		}
	}

	if preloads := opts.preloads; len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(preload)
		}
	}

	if !opts.skipPagination {
		if opts.page > 0 {
			db = db.Offset(int(opts.page * opts.size))
		}

		if opts.size > 0 {
			db = db.Limit(int(opts.size))
		}
	}

	if len(opts.filter) > 0 {
		for field, value := range opts.filter {
			if value.Raw != nil {
				db = db.Where(value.Raw)
			} else {
				db = db.Where(fmt.Sprintf("%s %s ?", field, value.Operator), value.Value)
			}
		}
	}

	if sort := opts.sort; sort != "" {
		db = db.Order(sort)
	}

	return db
}
