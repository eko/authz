package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/database/model"
	"github.com/eko/authz/backend/internal/event"
	"gorm.io/gorm"
)

const (
	WildcardValue = "*"
)

type Manager interface {
	CreateAction(identifier string) (*model.Action, error)
	CreateCompiledPolicy(compiledPolicy []*model.CompiledPolicy) error
	CreatePolicy(identifier string, resources []string, actions []string) (*model.Policy, error)
	CreatePrincipal(identifier string, roles []string) (*model.Principal, error)
	CreateResource(identifier string, kind string, value string) (*model.Resource, error)
	CreateRole(identifier string, policies []string) (*model.Role, error)
	GetActionRepository() *database.Repository[model.Action]
	GetCompiledPolicyRepository() *database.Repository[model.CompiledPolicy]
	GetPolicyRepository() *database.Repository[model.Policy]
	GetPrincipalRepository() *database.Repository[model.Principal]
	GetResourceRepository() *database.Repository[model.Resource]
	GetRoleRepository() *database.Repository[model.Role]
	IsAllowed(principalID string, resourceKind string, resourceValue string, actionID string) (bool, error)
	UpdatePolicy(identifier string, resources []string, actions []string) (*model.Policy, error)
	UpdatePrincipal(identifier string, roles []string) (*model.Principal, error)
	UpdateRole(identifier string, policies []string) (*model.Role, error)
}

type manager struct {
	dispatcher               event.Dispatcher
	actionRepository         *database.Repository[model.Action]
	compiledPolicyRepository *database.Repository[model.CompiledPolicy]
	policyRepository         *database.Repository[model.Policy]
	resourceRepository       *database.Repository[model.Resource]
	roleRepository           *database.Repository[model.Role]
	principalRepository      *database.Repository[model.Principal]
}

func New(
	dispatcher event.Dispatcher,
	actionRepository *database.Repository[model.Action],
	compiledPolicyRepository *database.Repository[model.CompiledPolicy],
	policyRepository *database.Repository[model.Policy],
	resourceRepository *database.Repository[model.Resource],
	roleRepository *database.Repository[model.Role],
	principalRepository *database.Repository[model.Principal],

) *manager {
	return &manager{
		dispatcher:               dispatcher,
		actionRepository:         actionRepository,
		compiledPolicyRepository: compiledPolicyRepository,
		policyRepository:         policyRepository,
		resourceRepository:       resourceRepository,
		roleRepository:           roleRepository,
		principalRepository:      principalRepository,
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

func (m *manager) CreateCompiledPolicy(compiledPolicy []*model.CompiledPolicy) error {
	if err := m.compiledPolicyRepository.Create(compiledPolicy...); err != nil {
		return fmt.Errorf("unable to create compiled policies: %v", err)
	}

	return nil
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

	if err := m.dispatcher.Dispatch(event.EventTypePolicy, policy.ID); err != nil {
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
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

	if err := m.dispatcher.Dispatch(event.EventTypePolicy, policy.ID); err != nil {
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
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
		value = WildcardValue
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

func (m *manager) CreatePrincipal(identifier string, roles []string) (*model.Principal, error) {
	exists, err := m.principalRepository.Get(identifier)
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

	principal := &model.Principal{
		ID:    identifier,
		Roles: roleObjects,
	}

	if err := m.principalRepository.Create(principal); err != nil {
		return nil, fmt.Errorf("unable to create principal: %v", err)
	}

	return principal, nil
}

func (m *manager) UpdatePrincipal(identifier string, roles []string) (*model.Principal, error) {
	principal, err := m.principalRepository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

	if err := m.principalRepository.Update(principal); err != nil {
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

func (m *manager) GetCompiledPolicyRepository() *database.Repository[model.CompiledPolicy] {
	return m.compiledPolicyRepository
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

func (m *manager) IsAllowed(principalID string, resourceKind string, resourceValue string, actionID string) (bool, error) {
	principal, err := m.principalRepository.Get(principalID, database.WithPreloads("Roles.Policies"))
	if err != nil {
		return false, fmt.Errorf("unable to retrieve principal: %v", err)
	}

	var policyIDs = make([]string, 0)
	for _, role := range principal.Roles {
		for _, policy := range role.Policies {
			policyIDs = append(policyIDs, policy.ID)
		}
	}

	return m.isPolicyAllowed(policyIDs, resourceKind, resourceValue, actionID)
}

func (m *manager) isPolicyAllowed(policyIDs []string, resourceKind string, resourceValue string, actionID string) (bool, error) {
	fields := map[string]database.FieldValue{
		"policy_id":      {Operator: "IN", Value: policyIDs},
		"resource_kind":  {Operator: "=", Value: resourceKind},
		"resource_value": {Operator: "=", Value: resourceValue},
		"action_id":      {Operator: "=", Value: actionID},
	}

	_, err := m.compiledPolicyRepository.GetByFields(fields)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if resourceValue != WildcardValue {
			return m.isPolicyAllowed(policyIDs, resourceKind, WildcardValue, actionID)
		}

		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("unable to retrieve compiled policies: %v", err)
	}

	return true, nil
}
