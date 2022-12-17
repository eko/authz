package database

import (
	"github.com/eko/authz/backend/internal/database/model"
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

// ResourceRepository struct that allows contacting the database using Gorm.
type ResourceRepository struct {
	*Repository[model.Resource]
}

// NewResourceRepository initializes a new repository.
func NewResourceRepository(repository *Repository[model.Resource]) *ResourceRepository {
	return &ResourceRepository{
		repository,
	}
}

type MatchingAttributeResourcePrincipal struct {
	PrincipalID   string
	ResourceKind  string
	ResourceValue string
}

func (r *ResourceRepository) FindMatchingAttributesWithPrincipals(
	resourceAttribute string,
	principalAttribute string,
	options ...ResourceQueryOption,
) ([]*MatchingAttributeResourcePrincipal, error) {
	matches := []*MatchingAttributeResourcePrincipal{}

	tx := applyResourceOptions(r.db, options)

	err := tx.
		Select("authz_principals_attributes.principal_id AS principal_id, authz_resources.kind AS resource_kind, authz_resources.value AS resource_value").
		Model(&model.Resource{}).
		Joins("LEFT JOIN authz_resources_attributes ON authz_resources.id = authz_resources_attributes.resource_id").
		Joins("LEFT JOIN authz_attributes ON authz_resources_attributes.attribute_id = authz_attributes.id").
		Joins("LEFT JOIN authz_principals_attributes ON authz_attributes.id = authz_principals_attributes.attribute_id").
		Where("authz_attributes.key = ? OR authz_attributes.key = ?", resourceAttribute, principalAttribute).
		Where("authz_principals_attributes.principal_id IS NOT NULL").
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
