package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/database/model"
	"gorm.io/gorm"
)

type Manager interface {
	CreateAction(identifier string) (*model.Action, error)
	CreatePolicy(identifier string, resources []string, actions []string) (*model.Policy, error)
	UpdatePolicy(identifier string, resources []string, actions []string) (*model.Policy, error)
	CreateResource(identifier string, kind string, value string) (*model.Resource, error)
	CreateRole(identifier string, policies []string) (*model.Role, error)
	UpdateRole(identifier string, policies []string) (*model.Role, error)
	CreatePrincipal(identifier string) (*model.Principal, error)
	GetActionRepository() *database.Repository[model.Action]
	GetPolicyRepository() *database.Repository[model.Policy]
	GetResourceRepository() *database.Repository[model.Resource]
	GetRoleRepository() *database.Repository[model.Role]
	GetPrincipalRepository() *database.Repository[model.Principal]
}

type manager struct {
	actionRepository    *database.Repository[model.Action]
	policyRepository    *database.Repository[model.Policy]
	resourceRepository  *database.Repository[model.Resource]
	roleRepository      *database.Repository[model.Role]
	principalRepository *database.Repository[model.Principal]
}

func New(
	actionRepository *database.Repository[model.Action],
	policyRepository *database.Repository[model.Policy],
	resourceRepository *database.Repository[model.Resource],
	roleRepository *database.Repository[model.Role],
	principalRepository *database.Repository[model.Principal],

) *manager {
	return &manager{
		actionRepository:    actionRepository,
		policyRepository:    policyRepository,
		resourceRepository:  resourceRepository,
		roleRepository:      roleRepository,
		principalRepository: principalRepository,
	}
}

func (m *manager) CreateAction(identifier string) (*model.Action, error) {
	exists, err := m.actionRepository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing action: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("an action already exists with identifier %q", identifier)
	}

	action := &model.Action{
		ID: identifier,
	}

	if err := m.actionRepository.Create(action); err != nil {
		return nil, fmt.Errorf("unable to create action: %v", err)
	}

	return action, nil
}

func (m *manager) CreatePolicy(identifier string, resources []string, actions []string) (*model.Policy, error) {
	exists, err := m.policyRepository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing policy: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a policy already exists with identifier %q", identifier)
	}

	policy := &model.Policy{}
	if err := m.attachToPolicy(policy, identifier, resources, actions); err != nil {
		return nil, err
	}

	if err := m.policyRepository.Create(policy); err != nil {
		return nil, fmt.Errorf("unable to create policy: %v", err)
	}

	return policy, nil
}

func (m *manager) UpdatePolicy(identifier string, resources []string, actions []string) (*model.Policy, error) {
	policy, err := m.policyRepository.Get(identifier)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve policy: %v", err)
	}

	if err := m.attachToPolicy(policy, policy.ID, resources, actions); err != nil {
		return nil, err
	}

	if err := m.policyRepository.Update(policy); err != nil {
		return nil, fmt.Errorf("unable to update policy: %v", err)
	}

	return policy, nil
}

func (m *manager) attachToPolicy(policy *model.Policy, identifier string, resources []string, actions []string) error {
	var resourceObjects = []*model.Resource{}

	for _, resource := range resources {
		resourceObject, err := m.resourceRepository.Get(resource)
		if err != nil {
			return fmt.Errorf("unable to retrieve resource %v: %v", resource, err)
		}

		resourceObjects = append(resourceObjects, resourceObject)
	}

	var actionObjects = []*model.Action{}

	for _, action := range actions {
		actionObject, err := m.actionRepository.Get(action)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			actionObject, err = m.CreateAction(action)
			if err != nil {
				return fmt.Errorf("unable to create action %q: %v", action, err)
			}
		} else if err != nil {
			return fmt.Errorf("unable to retrieve action %q: %v", action, err)
		}

		actionObjects = append(actionObjects, actionObject)
	}

	policy.ID = identifier
	policy.Resources = resourceObjects
	policy.Actions = actionObjects

	return nil
}

func (m *manager) CreateResource(identifier string, kind string, value string) (*model.Resource, error) {
	if value == "" {
		value = "*"
	}

	exists, err := m.resourceRepository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing resource: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a resource already exists with kind %q and value %q", kind, value)
	}

	resource := &model.Resource{
		ID:    identifier,
		Kind:  kind,
		Value: value,
	}

	if err := m.resourceRepository.Create(resource); err != nil {
		return nil, fmt.Errorf("unable to create resource: %v", err)
	}

	return resource, nil
}

func (m *manager) CreatePrincipal(identifier string) (*model.Principal, error) {
	exists, err := m.principalRepository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing principal: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a principal already exists with identifier %q", identifier)
	}

	principal := &model.Principal{
		ID: identifier,
	}

	if err := m.principalRepository.Create(principal); err != nil {
		return nil, fmt.Errorf("unable to create principal: %v", err)
	}

	return principal, nil
}

func (m *manager) CreateRole(identifier string, policies []string) (*model.Role, error) {
	exists, err := m.roleRepository.Get(identifier)
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

	if err := m.roleRepository.Create(role); err != nil {
		return nil, fmt.Errorf("unable to create role: %v", err)
	}

	return role, nil
}

func (m *manager) UpdateRole(identifier string, policies []string) (*model.Role, error) {
	role, err := m.roleRepository.Get(identifier)
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

	if err := m.roleRepository.Update(role); err != nil {
		return nil, fmt.Errorf("unable to update role: %v", err)
	}

	return role, nil
}

func (m *manager) GetActionRepository() *database.Repository[model.Action] {
	return m.actionRepository
}

func (m *manager) GetPolicyRepository() *database.Repository[model.Policy] {
	return m.policyRepository
}

func (m *manager) GetResourceRepository() *database.Repository[model.Resource] {
	return m.resourceRepository
}

func (m *manager) GetRoleRepository() *database.Repository[model.Role] {
	return m.roleRepository
}

func (m *manager) GetPrincipalRepository() *database.Repository[model.Principal] {
	return m.principalRepository
}
