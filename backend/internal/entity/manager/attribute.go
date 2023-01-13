package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"gorm.io/gorm"
)

type AttributeRepository repository.Base[model.Attribute]

type Attribute interface {
	MapToSlice(attributes map[string]any) ([]*model.Attribute, error)
	GetRepository() AttributeRepository
}

type attributeManager struct {
	repository AttributeRepository
}

// NewAttribute initializes a new attribute manager.
func NewAttribute(repository AttributeRepository) Attribute {
	return &attributeManager{
		repository: repository,
	}
}

func (m *attributeManager) GetRepository() AttributeRepository {
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
			"key_name": {Operator: "=", Value: attributeKey},
			"value":    {Operator: "=", Value: value},
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
