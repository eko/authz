package entity

import (
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func FxModule() fx.Option {
	return fx.Module("entity",
		fx.Provide(
			manager.NewAction,
			manager.NewAttribute,
			manager.NewAudit,
			manager.NewClient,
			manager.NewCompiledPolicy,
			manager.NewPolicy,
			manager.NewPrincipal,
			manager.NewResource,
			manager.NewRole,
			manager.NewStats,
			manager.NewUser,

			// Action
			func(db *gorm.DB) repository.Base[model.Action] {
				return repository.New[model.Action](db)
			},

			func(repository repository.Base[model.Action]) manager.ActionRepository {
				return repository
			},

			// Attribute
			func(db *gorm.DB) repository.Base[model.Attribute] {
				return repository.New[model.Attribute](db)
			},

			func(repository repository.Base[model.Attribute]) manager.AttributeRepository {
				return repository
			},

			// Audit
			func(db *gorm.DB) repository.Base[model.Audit] {
				return repository.New[model.Audit](db)
			},

			func(repository repository.Base[model.Audit]) manager.AuditRepository {
				return repository
			},

			// Client
			func(db *gorm.DB) repository.Base[model.Client] {
				return repository.New[model.Client](db)
			},

			func(repository repository.Base[model.Client]) manager.ClientRepository {
				return repository
			},

			// CompiledPolicy
			func(db *gorm.DB) repository.Base[model.CompiledPolicy] {
				return repository.New[model.CompiledPolicy](db)
			},

			func(repository repository.Base[model.CompiledPolicy]) manager.CompiledPolicyRepository {
				return repository
			},

			// Policy
			func(db *gorm.DB) repository.Base[model.Policy] {
				return repository.New[model.Policy](db)
			},

			func(repository repository.Base[model.Policy]) manager.PolicyRepository {
				return repository
			},

			// Principal
			func(db *gorm.DB) repository.Base[model.Principal] {
				return repository.New[model.Principal](db)
			},

			func(base repository.Base[model.Principal]) repository.Principal {
				return repository.NewPrincipal(base)
			},

			// Resource
			func(db *gorm.DB) repository.Base[model.Resource] {
				return repository.New[model.Resource](db)
			},

			func(base repository.Base[model.Resource]) repository.Resource {
				return repository.NewResource(base)
			},

			// Role
			func(db *gorm.DB) repository.Base[model.Role] {
				return repository.New[model.Role](db)
			},

			func(repository repository.Base[model.Role]) manager.RoleRepository {
				return repository
			},

			// Stats
			func(db *gorm.DB) repository.Base[model.Stats] {
				return repository.New[model.Stats](db)
			},

			func(repository repository.Base[model.Stats]) manager.StatsRepository {
				return repository
			},

			// User
			func(db *gorm.DB) repository.Base[model.User] {
				return repository.New[model.User](db)
			},

			func(repository repository.Base[model.User]) manager.UserRepository {
				return repository
			},
		),
	)
}
