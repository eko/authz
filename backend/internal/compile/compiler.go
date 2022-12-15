package compile

import (
	"fmt"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/database/model"
	"github.com/eko/authz/backend/internal/helper/time"
	"github.com/eko/authz/backend/internal/manager"
)

type Compiler interface {
	CompilePolicy(identifier string) error
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

	if err := c.manager.CreateCompiledPolicy(compiled); err != nil {
		return err
	}

	return c.manager.GetCompiledPolicyRepository().DeleteByFields(map[string]database.FieldValue{
		"policy_id": {Operator: "=", Value: policy.ID},
		"version":   {Operator: "<", Value: version},
	})
}
