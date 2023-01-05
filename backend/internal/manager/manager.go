package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/attribute"
	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/database/model"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/helper/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

const (
	WildcardValue = "*"
)

type Manager interface {
	CreateAction(identifier string) (*model.Action, error)
	CreateClient(name string, domain string) (*model.Client, error)
	CreateCompiledPolicy(compiledPolicy []*model.CompiledPolicy) error
	CreatePolicy(identifier string, resources []string, actions []string, attributeRules []string) (*model.Policy, error)
	CreatePrincipal(identifier string, roles []string, attributes map[string]any) (*model.Principal, error)
	CreateResource(identifier string, kind string, value string, attributes map[string]any) (*model.Resource, error)
	CreateRole(identifier string, policies []string) (*model.Role, error)
	CreateUser(username string, password string) (*model.User, error)
	GetActionRepository() *database.Repository[model.Action]
	GetClientRepository() *database.Repository[model.Client]
	GetCompiledPolicyRepository() *database.Repository[model.CompiledPolicy]
	GetPolicyRepository() *database.Repository[model.Policy]
	GetPrincipalRepository() *database.Repository[model.Principal]
	GetResourceRepository() *database.ResourceRepository
	GetRoleRepository() *database.Repository[model.Role]
	GetUserRepository() *database.Repository[model.User]
	IsAllowed(principalID string, resourceKind string, resourceValue string, actionID string) (bool, error)
	UpdatePolicy(identifier string, resources []string, actions []string, attributeRules []string) (*model.Policy, error)
	UpdatePrincipal(identifier string, roles []string, attributes map[string]any) (*model.Principal, error)
	UpdateRole(identifier string, policies []string) (*model.Role, error)
}

type manager struct {
	logger                   *slog.Logger
	db                       *gorm.DB
	transactionManager       database.TransactionManager
	dispatcher               event.Dispatcher
	tokenGenerator           token.Generator
	actionRepository         *database.Repository[model.Action]
	attributeRepository      *database.Repository[model.Attribute]
	clientRepository         *database.Repository[model.Client]
	compiledPolicyRepository *database.Repository[model.CompiledPolicy]
	policyRepository         *database.Repository[model.Policy]
	resourceRepository       *database.ResourceRepository
	roleRepository           *database.Repository[model.Role]
	principalRepository      *database.Repository[model.Principal]
	userRepository           *database.Repository[model.User]
}

func New(
	logger *slog.Logger,
	db *gorm.DB,
	transactionManager database.TransactionManager,
	dispatcher event.Dispatcher,
	tokenGenerator token.Generator,
	actionRepository *database.Repository[model.Action],
	attributeRepository *database.Repository[model.Attribute],
	clientRepository *database.Repository[model.Client],
	compiledPolicyRepository *database.Repository[model.CompiledPolicy],
	policyRepository *database.Repository[model.Policy],
	resourceRepository *database.ResourceRepository,
	roleRepository *database.Repository[model.Role],
	principalRepository *database.Repository[model.Principal],
	userRepository *database.Repository[model.User],
) *manager {
	return &manager{
		logger:                   logger,
		db:                       db,
		transactionManager:       transactionManager,
		dispatcher:               dispatcher,
		tokenGenerator:           tokenGenerator,
		actionRepository:         actionRepository,
		attributeRepository:      attributeRepository,
		clientRepository:         clientRepository,
		compiledPolicyRepository: compiledPolicyRepository,
		policyRepository:         policyRepository,
		resourceRepository:       resourceRepository,
		roleRepository:           roleRepository,
		principalRepository:      principalRepository,
		userRepository:           userRepository,
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

func (m *manager) CreateClient(name string, domain string) (*model.Client, error) {
	exists, err := m.clientRepository.GetByFields(map[string]database.FieldValue{
		"name": {Operator: "=", Value: name},
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing client: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a client already exists with name %q", name)
	}

	clientID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate client identifier: %v", err)
	}

	secret, err := m.tokenGenerator.Generate(48)
	if err != nil {
		return nil, fmt.Errorf("unable to generate client secret: %v", err)
	}

	client := &model.Client{
		ID:     clientID.String(),
		Secret: secret,
		Domain: domain,
		Name:   name,
	}

	if err := m.clientRepository.Create(client); err != nil {
		return nil, fmt.Errorf("unable to create client: %v", err)
	}

	return client, nil
}

func (m *manager) CreateCompiledPolicy(compiledPolicy []*model.CompiledPolicy) error {
	if err := m.compiledPolicyRepository.Create(compiledPolicy...); err != nil {
		return fmt.Errorf("unable to create compiled policies: %v", err)
	}

	return nil
}

func (m *manager) CreatePolicy(identifier string, resources []string, actions []string, attributeRules []string) (*model.Policy, error) {
	exists, err := m.policyRepository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing policy: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a policy already exists with identifier %q", identifier)
	}

	policy := &model.Policy{}
	if err := m.attachToPolicy(policy, identifier, resources, actions, attributeRules); err != nil {
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

func (m *manager) UpdatePolicy(identifier string, resources []string, actions []string, attributeRules []string) (*model.Policy, error) {
	policy, err := m.policyRepository.Get(identifier)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve policy: %v", err)
	}

	if err := m.attachToPolicy(policy, policy.ID, resources, actions, attributeRules); err != nil {
		return nil, err
	}

	transaction := m.transactionManager.New()
	defer func() { _ = transaction.Commit() }()

	policyRepository := m.policyRepository.WithTransaction(transaction)

	if err := policyRepository.UpdateAssociation(policy, "Resources", policy.Resources); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to update policy resources association: %v", err)
	}

	if err := policyRepository.UpdateAssociation(policy, "Actions", policy.Actions); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to update policy actions association: %v", err)
	}

	if err := policyRepository.Update(policy); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to update policy: %v", err)
	}

	if err := m.dispatcher.Dispatch(event.EventTypePolicy, policy.ID); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
	}

	return policy, nil
}

func (m *manager) attachToPolicy(
	policy *model.Policy,
	identifier string,
	resources []string,
	actions []string,
	attributeRules []string,
) error {
	for _, attributeRule := range attributeRules {
		if _, err := attribute.ConvertStringToRuleOperator(attributeRule); err != nil {
			return fmt.Errorf("unable to convert attribute rule %q to rule operator: %v", attributeRule, err)
		}
	}

	var resourceObjects = []*model.Resource{}

	for _, resource := range resources {
		resourceObject, err := m.resourceRepository.Get(resource)
		kind, value := ResourceSplit(resource)

		if errors.Is(err, gorm.ErrRecordNotFound) && value == WildcardValue {
			resourceObject, err = m.CreateResource(resource, kind, value, map[string]any{})
			if err != nil {
				return fmt.Errorf("unable to create wildcard resource %v: %v", resource, err)
			}
		} else if err != nil {
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
	policy.AttributeRules = attributeRules

	return nil
}

func (m *manager) CreateResource(identifier string, kind string, value string, attributes map[string]any) (*model.Resource, error) {
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

	attributeObjects, err := m.attributesMapToSlice(attributes)
	if err != nil {
		return nil, fmt.Errorf("unable to convert attributes to slice: %v", err)
	}

	resource := &model.Resource{
		ID:         identifier,
		Kind:       kind,
		Value:      value,
		Attributes: attributeObjects,
	}

	if err := m.resourceRepository.Create(resource); err != nil {
		return nil, fmt.Errorf("unable to create resource: %v", err)
	}

	if err := m.dispatcher.Dispatch(event.EventTypeResource, resource.ID); err != nil {
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
	}

	return resource, nil
}

func (m *manager) CreatePrincipal(identifier string, roles []string, attributes map[string]any) (*model.Principal, error) {
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

	attributeObjects, err := m.attributesMapToSlice(attributes)
	if err != nil {
		return nil, fmt.Errorf("unable to convert attributes to slice: %v", err)
	}

	principal := &model.Principal{
		ID:         identifier,
		Roles:      roleObjects,
		Attributes: attributeObjects,
	}

	if err := m.principalRepository.Create(principal); err != nil {
		return nil, fmt.Errorf("unable to create principal: %v", err)
	}

	if err := m.dispatcher.Dispatch(event.EventTypePrincipal, principal.ID); err != nil {
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
	}

	return principal, nil
}

func (m *manager) UpdatePrincipal(identifier string, roles []string, attributes map[string]any) (*model.Principal, error) {
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

	attributeObjects, err := m.attributesMapToSlice(attributes)
	if err != nil {
		return nil, fmt.Errorf("unable to convert attributes to slice: %v", err)
	}

	principal.Attributes = attributeObjects

	if err := m.principalRepository.Update(principal); err != nil {
		return nil, fmt.Errorf("unable to create principal: %v", err)
	}

	if err := m.dispatcher.Dispatch(event.EventTypePrincipal, principal.ID); err != nil {
		return nil, fmt.Errorf("unable to dispatch event: %v", err)
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

func (m *manager) CreateUser(username string, password string) (*model.User, error) {
	exists, err := m.userRepository.GetByFields(map[string]database.FieldValue{
		"username": {Operator: "=", Value: username},
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing user: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a user already exists with username %q", username)
	}

	if password == "" {
		password, err = m.tokenGenerator.Generate(10)
		if err != nil {
			return nil, fmt.Errorf("unable to generate a random password: %v", err)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	transaction := m.transactionManager.New()
	defer func() { _ = transaction.Commit() }()

	user := &model.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if err := m.userRepository.WithTransaction(transaction).Create(user); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to create user: %v", err)
	}

	if err := m.principalRepository.WithTransaction(transaction).Create(&model.Principal{
		ID: fmt.Sprintf("%s-%s", configs.ApplicationName, user.Username),
	}); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to create user: %v", err)
	}

	return user, nil
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

	transaction := m.transactionManager.New()
	defer func() { _ = transaction.Commit() }()

	roleRepository := m.roleRepository.WithTransaction(transaction)

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

func (m *manager) GetActionRepository() *database.Repository[model.Action] {
	return m.actionRepository
}

func (m *manager) GetClientRepository() *database.Repository[model.Client] {
	return m.clientRepository
}

func (m *manager) GetCompiledPolicyRepository() *database.Repository[model.CompiledPolicy] {
	return m.compiledPolicyRepository
}

func (m *manager) GetPolicyRepository() *database.Repository[model.Policy] {
	return m.policyRepository
}

func (m *manager) GetResourceRepository() *database.ResourceRepository {
	return m.resourceRepository
}

func (m *manager) GetRoleRepository() *database.Repository[model.Role] {
	return m.roleRepository
}

func (m *manager) GetPrincipalRepository() *database.Repository[model.Principal] {
	return m.principalRepository
}

func (m *manager) GetUserRepository() *database.Repository[model.User] {
	return m.userRepository
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

	isAllowed, err := m.isPolicyAllowed(policyIDs, resourceKind, resourceValue, actionID)
	if err != nil {
		return false, err
	}

	if !isAllowed {
		isAllowed, err = m.isPrincipalAllowed(principalID, resourceKind, resourceValue, actionID)
		if err != nil {
			return false, err
		}
	}

	m.logger.Debug(
		"Call to IsAllowed method",
		slog.String("principal_id", principalID),
		slog.String("resource_kind", resourceKind),
		slog.String("resource_value", resourceValue),
		slog.String("action_id", actionID),
		slog.Bool("result", isAllowed),
	)

	return isAllowed, err
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

func (m *manager) isPrincipalAllowed(principalID string, resourceKind string, resourceValue string, actionID string) (bool, error) {
	fields := map[string]database.FieldValue{
		"principal_id":   {Operator: "=", Value: principalID},
		"resource_kind":  {Operator: "=", Value: resourceKind},
		"resource_value": {Operator: "=", Value: resourceValue},
		"action_id":      {Operator: "=", Value: actionID},
	}

	_, err := m.compiledPolicyRepository.GetByFields(fields)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if resourceValue != WildcardValue {
			return m.isPrincipalAllowed(principalID, resourceKind, WildcardValue, actionID)
		}

		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("unable to retrieve compiled policies: %v", err)
	}

	return true, nil
}

func (m *manager) attributesMapToSlice(attributes map[string]any) ([]*model.Attribute, error) {
	var attributeObjects = make([]*model.Attribute, 0)

	for attributeKey, attributeValue := range attributes {
		value, err := CastAnyToString(attributeValue)
		if err != nil {
			return nil, fmt.Errorf("unable to cast attribute value to string: %v", err)
		}

		attribute, err := m.attributeRepository.GetByFields(map[string]database.FieldValue{
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
