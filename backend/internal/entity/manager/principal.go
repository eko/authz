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

type Principal interface {
	Create(identifier string, roles []string, attributes map[string]any) (*model.Principal, error)
	Delete(identifier string) error
	GetRepository() repository.Base[model.Principal]
	Update(identifier string, roles []string, attributes map[string]any) (*model.Principal, error)
}

type principalManager struct {
	repository         repository.Base[model.Principal]
	roleRepository     repository.Base[model.Role]
	attributeManager   Attribute
	transactionManager database.TransactionManager
	dispatcher         event.Dispatcher
}

// NewPrincipal initializes a new principal manager.
func NewPrincipal(
	repository repository.Base[model.Principal],
	roleRepository repository.Base[model.Role],
	attributeManager Attribute,
	transactionManager database.TransactionManager,
	dispatcher event.Dispatcher,
) Principal {
	return &principalManager{
		repository:         repository,
		roleRepository:     roleRepository,
		attributeManager:   attributeManager,
		transactionManager: transactionManager,
		dispatcher:         dispatcher,
	}
}

func (m *principalManager) GetRepository() repository.Base[model.Principal] {
	return m.repository
}

func (m *principalManager) Create(identifier string, roles []string, attributes map[string]any) (*model.Principal, error) {
	exists, err := m.repository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing principal: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a principal already exists with identifier %q", identifier)
	}

	var roleObjects = []*model.Role{}

	for _, role := range roles {
		roleObject, err := m.roleRepository.Get(role)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve role %v: %v", role, err)
		}

		roleObjects = append(roleObjects, roleObject)
	}

	attributeObjects, err := m.attributeManager.MapToSlice(attributes)
	if err != nil {
		return nil, fmt.Errorf("unable to convert attributes to slice: %v", err)
	}

	principal := &model.Principal{
		ID:         identifier,
		Roles:      roleObjects,
		Attributes: attributeObjects,
	}

	if err := m.repository.Create(principal); err != nil {
		return nil, fmt.Errorf("unable to create principal: %v", err)
	}

	if err := m.dispatcher.Dispatch(event.EventTypePrincipal, principal.ID); err != nil {
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
	}

	return principal, nil
}

func (m *principalManager) Update(identifier string, roles []string, attributes map[string]any) (*model.Principal, error) {
	principal, err := m.repository.Get(identifier)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve principal: %v", err)
	}

	var roleObjects = []*model.Role{}

	for _, role := range roles {
		roleObject, err := m.roleRepository.Get(role)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve role %v: %v", role, err)
		}

		roleObjects = append(roleObjects, roleObject)
	}

	principal.Roles = roleObjects

	attributeObjects, err := m.attributeManager.MapToSlice(attributes)
	if err != nil {
		return nil, fmt.Errorf("unable to convert attributes to slice: %v", err)
	}

	principal.Attributes = attributeObjects

	transaction := m.transactionManager.New()
	defer func() { _ = transaction.Commit() }()

	principalRepository := m.repository.WithTransaction(transaction)

	if err := principalRepository.UpdateAssociation(principal, "Roles", principal.Roles); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to update principal roles association: %v", err)
	}

	if err := principalRepository.UpdateAssociation(principal, "Attributes", principal.Attributes); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to update principal attributes association: %v", err)
	}

	if err := principalRepository.Update(principal); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to create principal: %v", err)
	}

	if err := m.dispatcher.Dispatch(event.EventTypePrincipal, principal.ID); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
	}

	return principal, nil
}

func (m *principalManager) Delete(identifier string) error {
	principal, err := m.repository.Get(identifier)
	if err != nil {
		return fmt.Errorf("unable to retrieve principal: %v", err)
	}

	if principal.IsLocked {
		return errors.New("cannot be deleted because it is locked")
	}

	if err := m.repository.Delete(principal); err != nil {
		return fmt.Errorf("cannot delete principal: %v", err)
	}

	return nil
}
