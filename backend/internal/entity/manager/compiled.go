package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/event"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type CompiledPolicyRepository repository.Base[model.CompiledPolicy]

type CompiledPolicy interface {
	Create(compiledPolicy []*model.CompiledPolicy) error
	GetRepository() CompiledPolicyRepository
	IsAllowed(principalID string, resourceKind string, resourceValue string, actionID string) (bool, error)
}

type compiledPolicyManager struct {
	repository          CompiledPolicyRepository
	principalRepository repository.Base[model.Principal]
	logger              *slog.Logger
	dispatcher          event.Dispatcher
}

// NewCompiledPolicy initializes a new compiledPolicy manager.
func NewCompiledPolicy(
	repository CompiledPolicyRepository,
	principalRepository repository.Base[model.Principal],
	logger *slog.Logger,
	dispatcher event.Dispatcher,
) CompiledPolicy {
	return &compiledPolicyManager{
		repository:          repository,
		principalRepository: principalRepository,
		logger:              logger,
		dispatcher:          dispatcher,
	}
}

func (m *compiledPolicyManager) GetRepository() CompiledPolicyRepository {
	return m.repository
}

func (m *compiledPolicyManager) Create(compiledPolicy []*model.CompiledPolicy) error {
	if err := m.repository.Create(compiledPolicy...); err != nil {
		return fmt.Errorf("unable to create compiled policies: %v", err)
	}

	return nil
}

func (m *compiledPolicyManager) IsAllowed(principalID string, resourceKind string, resourceValue string, actionID string) (bool, error) {
	principal, err := m.principalRepository.Get(principalID, repository.WithPreloads("Roles.Policies"))
	if err != nil {
		return false, fmt.Errorf("unable to retrieve principal: %v", err)
	}

	var policyIDs = make([]string, 0)
	for _, role := range principal.Roles {
		for _, policy := range role.Policies {
			policyIDs = append(policyIDs, policy.ID)
		}
	}

	isAllowed, compiledPolicy, err := m.isPolicyAllowed(policyIDs, resourceKind, resourceValue, actionID)
	if err != nil {
		return false, err
	}

	if !isAllowed {
		isAllowed, compiledPolicy, err = m.isPrincipalAllowed(principalID, resourceKind, resourceValue, actionID)
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

	if err := m.dispatcher.Dispatch(event.EventTypeCheck, &event.CheckEvent{
		Principal:      principalID,
		ResourceKind:   resourceKind,
		ResourceValue:  resourceValue,
		Action:         actionID,
		IsAllowed:      isAllowed,
		CompiledPolicy: compiledPolicy,
	}); err != nil {
		m.logger.Error("unable to dispatch check event", err)
	}

	return isAllowed, nil
}

func (m *compiledPolicyManager) isPolicyAllowed(policyIDs []string, resourceKind string, resourceValue string, actionID string) (bool, *model.CompiledPolicy, error) {
	fields := map[string]repository.FieldValue{
		"policy_id":      {Operator: "IN", Value: policyIDs},
		"resource_kind":  {Operator: "=", Value: resourceKind},
		"resource_value": {Operator: "=", Value: resourceValue},
		"action_id":      {Operator: "=", Value: actionID},
	}

	commiledPolicy, err := m.repository.GetByFields(fields)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if resourceValue != WildcardValue {
			return m.isPolicyAllowed(policyIDs, resourceKind, WildcardValue, actionID)
		}

		return false, nil, nil
	} else if err != nil {
		return false, nil, fmt.Errorf("unable to retrieve compiled policies: %v", err)
	}

	return true, commiledPolicy, nil
}

func (m *compiledPolicyManager) isPrincipalAllowed(principalID string, resourceKind string, resourceValue string, actionID string) (bool, *model.CompiledPolicy, error) {
	fields := map[string]repository.FieldValue{
		"principal_id":   {Operator: "=", Value: principalID},
		"resource_kind":  {Operator: "=", Value: resourceKind},
		"resource_value": {Operator: "=", Value: resourceValue},
		"action_id":      {Operator: "=", Value: actionID},
	}

	compiledPolicy, err := m.repository.GetByFields(fields)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if resourceValue != WildcardValue {
			return m.isPrincipalAllowed(principalID, resourceKind, WildcardValue, actionID)
		}

		return false, nil, nil
	} else if err != nil {
		return false, nil, fmt.Errorf("unable to retrieve compiled policies: %v", err)
	}

	return true, compiledPolicy, nil
}
