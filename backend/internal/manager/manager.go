package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/database/model"
	"gorm.io/gorm"
)

type Manager interface {
	CreateAction(name string) (*model.Action, error)
	CreatePolicy(name string, resources []map[string]string, actions []string) (*model.Policy, error)
	UpdatePolicy(identifier int64, name string, resources []map[string]string, actions []string) (*model.Policy, error)
	CreateResource(kind string, value string) (*model.Resource, error)
	CreatePrincipal(value string) (*model.Principal, error)
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

func (m *manager) CreateAction(name string) (*model.Action, error) {
	exists, err := m.actionRepository.GetByField("name", name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing action: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("an action already exists with name %q", name)
	}

	action := &model.Action{
		Name: name,
	}

	if err := m.actionRepository.Create(action); err != nil {
		return nil, fmt.Errorf("unable to create action: %v", err)
	}

	return action, nil
}

func (m *manager) CreatePolicy(name string, resources []map[string]string, actions []string) (*model.Policy, error) {
	exists, err := m.policyRepository.GetByField("name", name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing policy: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a policy already exists with name %q", name)
	}

	policy := &model.Policy{}
	policy, err = m.attachToPolicy(policy, name, resources, actions)
	if err != nil {
		return nil, err
	}

	if err := m.policyRepository.Create(policy); err != nil {
		return nil, fmt.Errorf("unable to create policy: %v", err)
	}

	return policy, nil
}

func (m *manager) UpdatePolicy(identifier int64, name string, resources []map[string]string, actions []string) (*model.Policy, error) {
	policy, err := m.policyRepository.Get(identifier)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve policy: %v", err)
	}

	policy, err = m.attachToPolicy(policy, name, resources, actions)
	if err != nil {
		return nil, err
	}

	if err := m.policyRepository.Update(policy); err != nil {
		return nil, fmt.Errorf("unable to update policy: %v", err)
	}

	return policy, nil
}

func (m *manager) attachToPolicy(policy *model.Policy, name string, resources []map[string]string, actions []string) (*model.Policy, error) {
	var resourceObjects = []*model.Resource{}

	for _, resource := range resources {
		resourceObject, err := m.resourceRepository.GetByFields(resource)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve resource %v: %v", resource, err)
		}

		resourceObjects = append(resourceObjects, resourceObject)
	}

	var actionObjects = []*model.Action{}

	for _, action := range actions {
		actionObject, err := m.actionRepository.GetByField("name", action)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			actionObject, err = m.CreateAction(action)
			if err != nil {
				return nil, fmt.Errorf("unable to create action %q: %v", action, err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("unable to retrieve action %q: %v", action, err)
		}

		actionObjects = append(actionObjects, actionObject)
	}

	policy.Name = name
	policy.Resources = resourceObjects
	policy.Actions = actionObjects

	return policy, nil
}

func (m *manager) CreateResource(kind string, value string) (*model.Resource, error) {
	if value == "" {
		value = "*"
	}

	exists, err := m.resourceRepository.GetByFields(map[string]string{
		"kind":  kind,
		"value": value,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing resource: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a resource already exists with kind %q and value %q", kind, value)
	}

	resource := &model.Resource{
		Kind:  kind,
		Value: value,
	}

	if err := m.resourceRepository.Create(resource); err != nil {
		return nil, fmt.Errorf("unable to create resource: %v", err)
	}

	return resource, nil
}

func (m *manager) CreatePrincipal(value string) (*model.Principal, error) {
	exists, err := m.principalRepository.GetByField("value", value)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing principal: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a principal already exists with value %q", value)
	}

	principal := &model.Principal{
		Value: value,
	}

	if err := m.principalRepository.Create(principal); err != nil {
		return nil, fmt.Errorf("unable to create principal: %v", err)
	}

	return principal, nil
}

func (m *manager) CreateRole(role *model.Role) error {
	return m.roleRepository.Create(role)
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
