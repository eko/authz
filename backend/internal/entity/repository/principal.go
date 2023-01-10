package repository

import (
	"github.com/eko/authz/backend/internal/entity/model"
)

type Principal interface {
	Base[model.Principal]
	FindMatchingAttribute(principalAttribute string) ([]*PrincipalMatchingAttribute, error)
}

// besource struct that allows contacting the database using Gorm.
type principal struct {
	Base[model.Principal]
}

// NewPrincipal initializes a new principal repository.
func NewPrincipal(repository Base[model.Principal]) Principal {
	return &principal{
		repository,
	}
}

type PrincipalMatchingAttribute struct {
	PrincipalID    string
	AttributeValue string
}

func (r *principal) FindMatchingAttribute(principalAttribute string) ([]*PrincipalMatchingAttribute, error) {
	matches := []*PrincipalMatchingAttribute{}

	err := r.DB().
		Select("authz_principals.id AS principal_id, authz_attributes.value AS attribute_value").
		Model(&model.Principal{}).
		Joins("INNER JOIN authz_principals_attributes ON authz_principals_attributes.principal_id = authz_principals.id").
		Joins("INNER JOIN authz_attributes ON authz_principals_attributes.attribute_id = authz_attributes.id").
		Where("authz_attributes.key = ?", principalAttribute).
		Scan(&matches).Error
	if err != nil {
		return nil, err
	}

	return matches, nil
}
