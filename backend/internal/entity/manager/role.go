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

type RoleRepository repository.Base[model.Role]

type Role interface {
	Create(identifier string, policies []string) (*model.Role, error)
	Delete(identifier string) error
	GetRepository() RoleRepository
	Update(identifier string, policies []string) (*model.Role, error)
}

type roleManager struct {
	repository         RoleRepository
	policyRepository   PolicyRepository
	transactionManager database.TransactionManager
	dispatcher         event.Dispatcher
}

// NewRole initializes a new role manager.
func NewRole(
	repository RoleRepository,
	policyRepository PolicyRepository,
	transactionManager database.TransactionManager,
	dispatcher event.Dispatcher,
) Role {
	return &roleManager{
		repository:         repository,
		policyRepository:   policyRepository,
		transactionManager: transactionManager,
		dispatcher:         dispatcher,
	}
}

func (m *roleManager) GetRepository() RoleRepository {
	return m.repository
}

func (m *roleManager) Create(identifier string, policies []string) (*model.Role, error) {
	exists, err := m.repository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing role: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a role already exists with identifier %q", identifier)
	}

	if len(policies) == 0 {
		return nil, fmt.Errorf("you have to specify at least one policy")
	}

	var policyObjects = []*model.Policy{}

	for _, policy := range policies {
		policyObject, err := m.policyRepository.Get(policy)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve policy %v: %v", policy, err)
		}

		policyObjects = append(policyObjects, policyObject)
	}

	role := &model.Role{
		ID:       identifier,
		Policies: policyObjects,
	}

	if err := m.repository.Create(role); err != nil {
		return nil, fmt.Errorf("unable to create role: %v", err)
	}

	if err := m.dispatcher.Dispatch(event.EventTypeRole, &event.ItemEvent{
		Action: event.ItemActionCreate,
		Data:   role,
	}); err != nil {
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
	}

	return role, nil
}

func (m *roleManager) Delete(identifier string) error {
	role, err := m.repository.Get(identifier)
	if err != nil {
		return fmt.Errorf("cannot retrieve role: %v", err)
	}

	if err := m.repository.Delete(role); err != nil {
		return fmt.Errorf("cannot delete role: %v", err)
	}

	return nil
}

func (m *roleManager) Update(identifier string, policies []string) (*model.Role, error) {
	role, err := m.repository.Get(identifier)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve role: %v", err)
	}

	var policyObjects = []*model.Policy{}

	for _, policy := range policies {
		policyObject, err := m.policyRepository.Get(policy)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve policy %v: %v", policy, err)
		}

		policyObjects = append(policyObjects, policyObject)
	}

	role.Policies = policyObjects

	transaction := m.transactionManager.New()
	defer func() { _ = transaction.Commit() }()

	roleRepository := m.repository.WithTransaction(transaction)

	if err := roleRepository.UpdateAssociation(role, "Policies", role.Policies); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to update role policies association: %v", err)
	}

	if err := roleRepository.Update(role); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to update role: %v", err)
	}

	if err := m.dispatcher.Dispatch(event.EventTypeRole, &event.ItemEvent{
		Action: event.ItemActionUpdate,
		Data:   role,
	}); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
	}

	return role, nil
}
