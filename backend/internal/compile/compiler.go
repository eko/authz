package compile

import (
	"fmt"

	"github.com/eko/authz/backend/internal/attribute"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/helper/time"
)

type CompileOption func(*compileOptions)

type compileOptions struct {
	resources  []*model.Resource
	principals []*model.Principal
}

func WithResources(resources ...*model.Resource) CompileOption {
	return func(o *compileOptions) {
		o.resources = resources
	}
}

func WithPrincipals(principals ...*model.Principal) CompileOption {
	return func(o *compileOptions) {
		o.principals = principals
	}
}

type Compiler interface {
	CompilePolicy(policy *model.Policy) error
	CompilePrincipal(principal *model.Principal) error
	CompileResource(resource *model.Resource) error
}

type compiler struct {
	clock            time.Clock
	compiledManager  manager.CompiledPolicy
	policyManager    manager.Policy
	principalManager manager.Principal
	resourceManager  manager.Resource
}

func NewCompiler(
	clock time.Clock,
	compiledManager manager.CompiledPolicy,
	policyManager manager.Policy,
	principalManager manager.Principal,
	resourceManager manager.Resource,
) *compiler {
	return &compiler{
		clock:            clock,
		compiledManager:  compiledManager,
		policyManager:    policyManager,
		principalManager: principalManager,
		resourceManager:  resourceManager,
	}
}

func (c *compiler) CompilePolicy(policy *model.Policy) error {
	policy, err := c.policyManager.GetRepository().Get(
		policy.ID,
		repository.WithPreloads("Resources", "Actions"),
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve policy: %v", err)
	}

	// In case policy has attribute rules, just compile them.
	if len(policy.AttributeRules.Data()) > 0 {
		return c.compilePolicyAttributes(policy)
	}

	if len(policy.Resources) == 0 || len(policy.Actions) == 0 {
		// Nothing to update
		return nil
	}

	version := c.clock.Now().Unix()

	var compiled = make([]*model.CompiledPolicy, 0)
	for _, resource := range policy.Resources {
		for _, action := range policy.Actions {
			if len(compiled) == 100 {
				if err := c.compiledManager.Create(compiled); err != nil {
					return err
				}
				compiled = make([]*model.CompiledPolicy, 0)
			}

			compiled = append(compiled, &model.CompiledPolicy{
				PolicyID:      policy.ID,
				ResourceKind:  resource.Kind,
				ResourceValue: resource.Value,
				ActionID:      action.ID,
				Version:       version,
			})
		}
	}

	if len(compiled) == 0 {
		return nil
	}

	if err := c.compiledManager.Create(compiled); err != nil {
		return err
	}

	return c.compiledManager.GetRepository().DeleteByFields(map[string]repository.FieldValue{
		"policy_id": {Operator: "=", Value: policy.ID},
		"version":   {Operator: "<", Value: version},
	})
}

func (c *compiler) compilePolicyAttributes(policy *model.Policy) error {
	version := c.clock.Now().Unix()

	for _, attributeRuleStr := range policy.AttributeRules.Data() {
		attributeRule, err := attribute.ConvertStringToRuleOperator(attributeRuleStr)
		if err != nil {
			return fmt.Errorf("cannot convert attribute rule string to object: %v", err)
		}

		if attributeRule.Value != "" {
			if err := c.compilePolicyAttributesWithValue(policy, attributeRule, version); err != nil {
				return err
			}
		} else {
			if err := c.compilePolicyAttributesWithMatching(policy, attributeRule, version); err != nil {
				return err
			}
		}
	}

	return c.compiledManager.GetRepository().DeleteByFields(map[string]repository.FieldValue{
		"policy_id": {Operator: "=", Value: policy.ID},
		"version":   {Operator: "<", Value: version},
	})
}

func (c *compiler) compilePolicyAttributesWithValue(
	policy *model.Policy,
	attributeRule *attribute.Rule,
	version int64,
	options ...CompileOption,
) (err error) {
	opts := applyOptions(options)

	var compiled = make([]*model.CompiledPolicy, 0)

	if opts.resources == nil {
		opts.resources, err = c.retrieveResources(policy.Resources, attributeRule)
		if err != nil {
			return fmt.Errorf("cannot retrieve resources: %v", err)
		}
	}

	if opts.principals == nil {
		opts.principals, err = c.retrievePrincipals(attributeRule)
		if err != nil {
			return fmt.Errorf("cannot retrieve principals: %v", err)
		}
	}

	for _, resource := range opts.resources {
		if resource.Value == manager.WildcardValue {
			// Don't handle wildcard resources to compiled policies
			// in case of attribute rules.
			continue
		}

		for _, principal := range opts.principals {
			for _, action := range policy.Actions {
				if len(compiled) == 100 {
					if err := c.compiledManager.Create(compiled); err != nil {
						return err
					}
					compiled = make([]*model.CompiledPolicy, 0)
				}

				compiled = append(compiled, &model.CompiledPolicy{
					PolicyID:      policy.ID,
					PrincipalID:   principal.ID,
					ResourceKind:  resource.Kind,
					ResourceValue: resource.Value,
					ActionID:      action.ID,
					Version:       version,
				})
			}
		}
	}

	if len(compiled) == 0 {
		return nil
	}

	return c.compiledManager.Create(compiled)
}

func (c *compiler) compilePolicyAttributesWithMatching(
	policy *model.Policy,
	attributeRule *attribute.Rule,
	version int64,
	options ...CompileOption,
) (err error) {
	opts := applyOptions(options)

	var compiled = make([]*model.CompiledPolicy, 0)

	queryOptions := []repository.ResourceQueryOption{}

	if len(opts.resources) > 0 {
		var resourceIDs = make([]string, len(opts.resources))
		for index, resource := range opts.resources {
			resourceIDs[index] = resource.ID
		}

		queryOptions = append(queryOptions, repository.WithResourceIDs(resourceIDs))
	}

	resourcesMatches, err := c.resourceManager.GetRepository().FindMatchingAttribute(
		attributeRule.ResourceAttribute,
		queryOptions...,
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve resource and principals matches: %v", err)
	}

	principalsMatches, err := c.principalManager.GetRepository().FindMatchingAttribute(
		attributeRule.PrincipalAttribute,
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve resource and principals matches: %v", err)
	}

	for _, resourceMatch := range resourcesMatches {
		for _, principalMatch := range principalsMatches {
			if resourceMatch.AttributeValue != principalMatch.AttributeValue {
				continue
			}

			for _, action := range policy.Actions {
				if len(compiled) == 100 {
					if err := c.compiledManager.Create(compiled); err != nil {
						return err
					}
					compiled = make([]*model.CompiledPolicy, 0)
				}

				compiled = append(compiled, &model.CompiledPolicy{
					PolicyID:      policy.ID,
					PrincipalID:   principalMatch.PrincipalID,
					ResourceKind:  resourceMatch.ResourceKind,
					ResourceValue: resourceMatch.ResourceValue,
					ActionID:      action.ID,
					Version:       version,
				})
			}
		}
	}

	if len(compiled) == 0 {
		return nil
	}

	return c.compiledManager.Create(compiled)
}

func (c *compiler) retrieveResources(resources []*model.Resource, rule *attribute.Rule) ([]*model.Resource, error) {
	var result = make([]*model.Resource, 0)

	for _, resource := range resources {
		if resource.Value != manager.WildcardValue {
			result = append(result, resource)
			continue
		}

		var filters = map[string]repository.FieldValue{
			"authz_resources.kind": {Operator: "=", Value: resource.Kind},
			// Don't handle wildcard resources to compiled policies
			// in case of attribute rules.
			"authz_resources.value": {Operator: "<>", Value: manager.WildcardValue},
		}

		if rule.ResourceAttribute != "" && rule.Value != "" {
			filters["authz_attributes.key_name"] = repository.FieldValue{
				Operator: "=", Value: rule.ResourceAttribute,
			}
		}

		allResources, _, err := c.resourceManager.GetRepository().Find(
			repository.WithJoin(
				"LEFT JOIN authz_resources_attributes ON authz_resources.id = authz_resources_attributes.resource_id",
				"LEFT JOIN authz_attributes ON authz_resources_attributes.attribute_id = authz_attributes.id",
			),
			repository.WithFilter(filters),
			repository.WithPreloads("Attributes"),
		)
		if err != nil {
			return nil, err
		}

		matchingResources := []*model.Resource{}

		for _, resource := range allResources {
			if !rule.MatchResource(resource.Attributes) {
				continue
			}

			matchingResources = append(matchingResources, resource)
		}

		result = append(result, matchingResources...)
	}

	return result, nil
}

func (c *compiler) retrievePrincipals(rule *attribute.Rule) ([]*model.Principal, error) {
	var filters = map[string]repository.FieldValue{}

	if rule.PrincipalAttribute != "" && rule.Value != "" {
		filters["authz_attributes.key_name"] = repository.FieldValue{
			Operator: "=", Value: rule.PrincipalAttribute,
		}
	}

	allPrincipals, _, err := c.principalManager.GetRepository().Find(
		repository.WithJoin(
			"LEFT JOIN authz_principals_attributes ON authz_principals.id = authz_principals_attributes.principal_id",
			"LEFT JOIN authz_attributes ON authz_principals_attributes.attribute_id = authz_attributes.id",
		),
		repository.WithFilter(filters),
		repository.WithPreloads("Attributes"),
	)
	if err != nil {
		return nil, err
	}

	matchingPrincipals := []*model.Principal{}

	for _, principal := range allPrincipals {
		if !rule.MatchPrincipal(principal.Attributes) {
			continue
		}

		matchingPrincipals = append(matchingPrincipals, principal)
	}

	return matchingPrincipals, nil
}

func (c *compiler) CompilePrincipal(principal *model.Principal) error {
	principal, err := c.principalManager.GetRepository().Get(principal.ID)
	if err != nil {
		return fmt.Errorf("cannot retrieve principal: %v", err)
	}

	version := c.clock.Now().Unix()

	policies, _, err := c.policyManager.GetRepository().Find(
		repository.WithPreloads("Resources", "Actions"),
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve policies: %v", err)
	}

	for _, policy := range policies {
		for _, attributeRuleStr := range policy.AttributeRules.Data() {
			attributeRule, err := attribute.ConvertStringToRuleOperator(attributeRuleStr)
			if err != nil {
				return fmt.Errorf("cannot convert attribute rule string to object: %v", err)
			}

			if attributeRule.Value != "" {
				if err := c.compilePolicyAttributesWithValue(policy, attributeRule, version); err != nil {
					return err
				}
			} else {
				if err := c.compilePolicyAttributesWithMatching(policy, attributeRule, version); err != nil {
					return err
				}
			}
		}
	}

	return c.compiledManager.GetRepository().DeleteByFields(map[string]repository.FieldValue{
		"principal_id": {Operator: "=", Value: principal.ID},
		"version":      {Operator: "<", Value: version},
	})
}

func (c *compiler) CompileResource(resource *model.Resource) error {
	resource, err := c.resourceManager.GetRepository().Get(resource.ID)
	if err != nil {
		return fmt.Errorf("cannot retrieve resource: %v", err)
	}

	version := c.clock.Now().Unix()

	policies, _, err := c.policyManager.GetRepository().Find(
		repository.WithPreloads("Resources", "Actions"),
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve policies: %v", err)
	}

	for _, policy := range policies {
		for _, attributeRuleStr := range policy.AttributeRules.Data() {
			attributeRule, err := attribute.ConvertStringToRuleOperator(attributeRuleStr)
			if err != nil {
				return fmt.Errorf("cannot convert attribute rule string to object: %v", err)
			}

			if attributeRule.Value != "" {
				if err := c.compilePolicyAttributesWithValue(policy, attributeRule, version, WithResources(resource)); err != nil {
					return err
				}
			} else {
				if err := c.compilePolicyAttributesWithMatching(policy, attributeRule, version, WithResources(resource)); err != nil {
					return err
				}
			}
		}
	}

	return c.compiledManager.GetRepository().DeleteByFields(map[string]repository.FieldValue{
		"resource_kind":  {Operator: "=", Value: resource.Kind},
		"resource_value": {Operator: "=", Value: resource.Value},
		"version":        {Operator: "<", Value: version},
	})
}

func applyOptions(options []CompileOption) *compileOptions {
	opts := &compileOptions{}

	for _, option := range options {
		option(opts)
	}

	return opts
}
