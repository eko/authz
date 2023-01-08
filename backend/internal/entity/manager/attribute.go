package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"gorm.io/gorm"
)

type Attribute interface {
	MapToSlice(attributes map[string]any) ([]*model.Attribute, error)
	GetRepository() repository.Base[model.Attribute]
}

type attributeManager struct {
	repository repository.Base[model.Attribute]
}

// NewAttribute initializes a new attribute manager.
func NewAttribute(repository repository.Base[model.Attribute]) Attribute {
	return &attributeManager{
		repository: repository,
	}
}

func (m *attributeManager) GetRepository() repository.Base[model.Attribute] {
	return m.repository
}

func (m *attributeManager) MapToSlice(attributes map[string]any) ([]*model.Attribute, error) {
	var attributeObjects = make([]*model.Attribute, 0)

	for attributeKey, attributeValue := range attributes {
		value, err := CastAnyToString(attributeValue)
		if err != nil {
			return nil, fmt.Errorf("unable to cast attribute value to string: %v", err)
		}

		attribute, err := m.repository.GetByFields(map[string]repository.FieldValue{
			"key":   {Operator: "=", Value: attributeKey},
			"value": {Operator: "=", Value: value},
		})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("unable to check for existing attribute: %v", err)
		}
		if attribute == nil {
			attribute = &model.Attribute{Key: attributeKey, Value: value}
		}

		attributeObjects = append(attributeObjects, attribute)
	}

	return attributeObjects, nil
}
