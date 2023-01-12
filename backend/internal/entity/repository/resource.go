package repository

import (
	"github.com/eko/authz/backend/internal/entity/model"
	"gorm.io/gorm"
)

type ResourceQueryOption func(*resourceQueryOptions)

type resourceQueryOptions struct {
	resourceIDs []string
}

func WithResourceIDs(resourceIDs []string) ResourceQueryOption {
	return func(o *resourceQueryOptions) {
		o.resourceIDs = resourceIDs
	}
}

type Resource interface {
	Base[model.Resource]
	FindMatchingAttribute(resourceAttribute string, options ...ResourceQueryOption) ([]*ResourceMatchingAttribute, error)
}

// besource struct that allows contacting the database using Gorm.
type resource struct {
	Base[model.Resource]
}

// NewResource initializes a new resource repository.
func NewResource(repository Base[model.Resource]) Resource {
	return &resource{
		repository,
	}
}

type ResourceMatchingAttribute struct {
	ResourceKind   string
	ResourceValue  string
	AttributeValue string
}

func (r *resource) FindMatchingAttribute(
	resourceAttribute string,
	options ...ResourceQueryOption,
) ([]*ResourceMatchingAttribute, error) {
	matches := []*ResourceMatchingAttribute{}

	tx := applyResourceOptions(r.DB(), options)

	err := tx.
		Select("authz_resources.kind AS resource_kind, authz_resources.value AS resource_value, authz_attributes.value AS attribute_value").
		Model(&model.Resource{}).
		Joins("INNER JOIN authz_resources_attributes ON authz_resources.id = authz_resources_attributes.resource_id").
		Joins("INNER JOIN authz_attributes ON authz_resources_attributes.attribute_id = authz_attributes.id").
		Where("authz_attributes.key_name = ?", resourceAttribute).
		Where("authz_resources.value <> ?", "*").
		Scan(&matches).Error
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func applyResourceOptions(tx *gorm.DB, options []ResourceQueryOption) *gorm.DB {
	opts := &resourceQueryOptions{}

	for _, opt := range options {
		opt(opts)
	}

	if len(opts.resourceIDs) > 0 {
		tx = tx.Where("authz_resources.id IN ?", opts.resourceIDs)
	}

	return tx
}
