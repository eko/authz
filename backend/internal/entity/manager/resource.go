package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/event"
	"gorm.io/gorm"
)

type Resource interface {
	Create(identifier string, kind string, value string, attributes map[string]any) (*model.Resource, error)
	Delete(identifier string) error
	GetRepository() repository.Resource
	Update(identifier string, kind string, value string, attributes map[string]any) (*model.Resource, error)
}

type resourceManager struct {
	repository         repository.Resource
	attributeManager   Attribute
	transactionManager database.TransactionManager
	dispatcher         event.Dispatcher
}

// NewResource initializes a new resource manager.
func NewResource(
	repository repository.Resource,
	attributeManager Attribute,
	transactionManager database.TransactionManager,
	dispatcher event.Dispatcher,
) Resource {
	return &resourceManager{
		repository:         repository,
		attributeManager:   attributeManager,
		transactionManager: transactionManager,
		dispatcher:         dispatcher,
	}
}

func (m *resourceManager) GetRepository() repository.Resource {
	return m.repository
}

func (m *resourceManager) Create(identifier string, kind string, value string, attributes map[string]any) (*model.Resource, error) {
	if value == "" {
		value = WildcardValue
	}

	exists, err := m.repository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing resource: %v", err)
	}

	existsKindValue, err := m.repository.GetByFields(map[string]repository.FieldValue{
		"kind":  {Operator: "=", Value: kind},
		"value": {Operator: "=", Value: value},
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing resource: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a resource already exists with id %q", identifier)
	}

	if existsKindValue != nil {
		return nil, fmt.Errorf("a resource already exists with kind %q and value %q", kind, value)
	}

	attributeObjects, err := m.attributeManager.MapToSlice(attributes)
	if err != nil {
		return nil, fmt.Errorf("unable to convert attributes to slice: %v", err)
	}

	resource := &model.Resource{
		ID:         identifier,
		Kind:       kind,
		Value:      value,
		Attributes: attributeObjects,
	}

	if err := m.repository.Create(resource); err != nil {
		return nil, fmt.Errorf("unable to create resource: %v", err)
	}

	if err := m.dispatcher.Dispatch(event.EventTypeResource, &event.ItemEvent{
		Action: event.ItemActionCreate,
		Data:   resource,
	}); err != nil {
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
	}

	return resource, nil
}

func (m *resourceManager) Delete(identifier string) error {
	resource, err := m.repository.Get(identifier)
	if err != nil {
		return fmt.Errorf("cannot retrieve resource: %v", err)
	}

	if err := m.repository.Delete(resource); err != nil {
		return fmt.Errorf("cannot delete resource: %v", err)
	}

	return nil
}

func (m *resourceManager) Update(identifier string, kind string, value string, attributes map[string]any) (*model.Resource, error) {
	resource, err := m.repository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to retrieve resource: %v", err)
	}

	attributeObjects, err := m.attributeManager.MapToSlice(attributes)
	if err != nil {
		return nil, fmt.Errorf("unable to convert attributes to slice: %v", err)
	}

	resource.Kind = kind
	resource.Value = value
	resource.Attributes = attributeObjects

	transaction := m.transactionManager.New()
	defer func() { _ = transaction.Commit() }()

	resourceRepository := m.repository.WithTransaction(transaction)

	if err := resourceRepository.UpdateAssociation(resource, "Attributes", resource.Attributes); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to update resource attributes association: %v", err)
	}

	if err := resourceRepository.Update(resource); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to update resource: %v", err)
	}

	if err := m.dispatcher.Dispatch(event.EventTypeResource, &event.ItemEvent{
		Action: event.ItemActionUpdate,
		Data:   resource,
	}); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
	}

	return resource, nil
}
