package compile

import (
	"fmt"

	"github.com/eko/authz/backend/internal/attribute"
	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/database/model"
	"github.com/eko/authz/backend/internal/helper/time"
	"github.com/eko/authz/backend/internal/manager"
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
	CompilePolicy(identifier string) error
	CompilePrincipal(identifier string) error
	CompileResource(identifier string) error
}

type compiler struct {
	clock   time.Clock
	manager manager.Manager
}

func NewCompiler(
	clock time.Clock,
	manager manager.Manager,
) *compiler {
	return &compiler{
		clock:   clock,
		manager: manager,
	}
}

func (c *compiler) CompilePolicy(identifier string) error {
	policy, err := c.manager.GetPolicyRepository().Get(
		identifier,
		database.WithPreloads("Resources", "Actions"),
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve policy: %v", err)
	}

	// In case policy has attribute rules, just compile them.
	if len(policy.AttributeRules) > 0 {
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
				if err := c.manager.CreateCompiledPolicy(compiled); err != nil {
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

	if err := c.manager.CreateCompiledPolicy(compiled); err != nil {
		return err
	}

	return c.manager.GetCompiledPolicyRepository().DeleteByFields(map[string]database.FieldValue{
		"policy_id": {Operator: "=", Value: policy.ID},
		"version":   {Operator: "<", Value: version},
	})
}

func (c *compiler) compilePolicyAttributes(policy *model.Policy) error {
	version := c.clock.Now().Unix()

	for _, attributeRuleStr := range policy.AttributeRules {
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

	return c.manager.GetCompiledPolicyRepository().DeleteByFields(map[string]database.FieldValue{
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
		for _, principal := range opts.principals {
			for _, action := range policy.Actions {
				if len(compiled) == 100 {
					if err := c.manager.CreateCompiledPolicy(compiled); err != nil {
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

	return c.manager.CreateCompiledPolicy(compiled)
}

func (c *compiler) compilePolicyAttributesWithMatching(
	policy *model.Policy,
	attributeRule *attribute.Rule,
	version int64,
	options ...CompileOption,
) (err error) {
	opts := applyOptions(options)

	var compiled = make([]*model.CompiledPolicy, 0)

	queryOptions := []database.ResourceQueryOption{}

	if len(opts.resources) > 0 {
		var resourceIDs = make([]string, len(opts.resources))
		for index, resource := range opts.resources {
			resourceIDs[index] = resource.ID
		}

		queryOptions = append(queryOptions, database.WithResourceIDs(resourceIDs))
	}

	matches, err := c.manager.GetResourceRepository().FindMatchingAttributesWithPrincipals(
		attributeRule.ResourceAttribute,
		attributeRule.PrincipalAttribute,
		queryOptions...,
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve resource and principals matches: %v", err)
	}

	for _, match := range matches {
		for _, action := range policy.Actions {
			if len(compiled) == 100 {
				if err := c.manager.CreateCompiledPolicy(compiled); err != nil {
					return err
				}
				compiled = make([]*model.CompiledPolicy, 0)
			}

			compiled = append(compiled, &model.CompiledPolicy{
				PolicyID:      policy.ID,
				PrincipalID:   match.PrincipalID,
				ResourceKind:  match.ResourceKind,
				ResourceValue: match.ResourceValue,
				ActionID:      action.ID,
				Version:       version,
			})
		}
	}

	if len(compiled) == 0 {
		return nil
	}

	return c.manager.CreateCompiledPolicy(compiled)
}

func (c *compiler) retrieveResources(resources []*model.Resource, rule *attribute.Rule) ([]*model.Resource, error) {
	var result = make([]*model.Resource, 0)

	for _, resource := range resources {
		if resource.Value != manager.WildcardValue {
			result = append(result, resource)
			continue
		}

		var filters = map[string]database.FieldValue{
			"kind": {Operator: "=", Value: resource.Kind},
		}

		if rule.ResourceAttribute != "" && rule.Value != "" {
			filters["authz_attributes.key"] = database.FieldValue{
				Operator: "=", Value: rule.ResourceAttribute,
			}

			switch rule.Operator {
			case attribute.RuleOperatorEqual:
				filters["authz_attributes.value"] = database.FieldValue{
					Operator: "=", Value: rule.Value,
				}
			case attribute.RuleOperatorNotEqual:
				filters["authz_attributes.value"] = database.FieldValue{
					Operator: "<>", Value: rule.Value,
				}
			}
		}

		allResources, _, err := c.manager.GetResourceRepository().Find(
			database.WithJoin(
				"INNER JOIN authz_resources_attributes ON authz_resources.id = authz_resources_attributes.resource_id",
				"INNER JOIN authz_attributes ON authz_resources_attributes.attribute_id = authz_attributes.id",
			),
			database.WithFilter(filters),
		)
		if err != nil {
			return nil, err
		}

		result = append(result, allResources...)
	}

	return result, nil
}

func (c *compiler) retrievePrincipals(rule *attribute.Rule) ([]*model.Principal, error) {
	var filters = map[string]database.FieldValue{}

	if rule.PrincipalAttribute != "" && rule.Value != "" {
		filters["authz_attributes.key"] = database.FieldValue{
			Operator: "=", Value: rule.PrincipalAttribute,
		}

		switch rule.Operator {
		case attribute.RuleOperatorEqual:
			filters["authz_attributes.value"] = database.FieldValue{
				Operator: "=", Value: rule.Value,
			}
		case attribute.RuleOperatorNotEqual:
			filters["authz_attributes.value"] = database.FieldValue{
				Operator: "<>", Value: rule.Value,
			}
		}
	}

	allPrincipals, _, err := c.manager.GetPrincipalRepository().Find(
		database.WithJoin(
			"INNER JOIN authz_principals_attributes ON authz_principals.id = authz_principals_attributes.principal_id",
			"INNER JOIN authz_attributes ON authz_principals_attributes.attribute_id = authz_attributes.id",
		),
		database.WithFilter(filters),
	)
	if err != nil {
		return nil, err
	}

	return allPrincipals, nil
}

func (c *compiler) CompilePrincipal(identifier string) error {
	principal, err := c.manager.GetPrincipalRepository().Get(identifier)
	if err != nil {
		return fmt.Errorf("cannot retrieve principal: %v", err)
	}

	version := c.clock.Now().Unix()

	policies, _, err := c.manager.GetPolicyRepository().Find(
		database.WithPreloads("Resources", "Actions"),
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve policies: %v", err)
	}

	for _, policy := range policies {
		for _, attributeRuleStr := range policy.AttributeRules {
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

	return c.manager.GetCompiledPolicyRepository().DeleteByFields(map[string]database.FieldValue{
		"principal_id": {Operator: "=", Value: principal.ID},
		"version":      {Operator: "<", Value: version},
	})
}

func (c *compiler) CompileResource(identifier string) error {
	resource, err := c.manager.GetResourceRepository().Get(identifier)
	if err != nil {
		return fmt.Errorf("cannot retrieve resource: %v", err)
	}

	version := c.clock.Now().Unix()

	policies, _, err := c.manager.GetPolicyRepository().Find(
		database.WithPreloads("Resources", "Actions"),
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve policies: %v", err)
	}

	for _, policy := range policies {
		for _, attributeRuleStr := range policy.AttributeRules {
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

	return c.manager.GetCompiledPolicyRepository().DeleteByFields(map[string]database.FieldValue{
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
