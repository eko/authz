package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"gorm.io/gorm"
)

type Role interface {
	Create(identifier string, policies []string) (*model.Role, error)
	GetRepository() repository.Base[model.Role]
	Update(identifier string, policies []string) (*model.Role, error)
}

type roleManager struct {
	repository         repository.Base[model.Role]
	policyRepository   repository.Base[model.Policy]
	transactionManager database.TransactionManager
}

// NewRole initializes a new role manager.
func NewRole(
	repository repository.Base[model.Role],
	policyRepository repository.Base[model.Policy],
	transactionManager database.TransactionManager,
) Role {
	return &roleManager{
		repository:         repository,
		policyRepository:   policyRepository,
		transactionManager: transactionManager,
	}
}

func (m *roleManager) GetRepository() repository.Base[model.Role] {
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

	return role, nil
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

	return role, nil
}
