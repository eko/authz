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
			manager.NewClient,
			manager.NewCompiledPolicy,
			manager.NewPolicy,
			manager.NewPrincipal,
			manager.NewResource,
			manager.NewRole,
			manager.NewStats,
			manager.NewUser,

			func(db *gorm.DB) repository.Base[model.Action] {
				return repository.New[model.Action](db)
			},

			func(db *gorm.DB) repository.Base[model.Attribute] {
				return repository.New[model.Attribute](db)
			},

			func(db *gorm.DB) repository.Base[model.Client] {
				return repository.New[model.Client](db)
			},

			func(db *gorm.DB) repository.Base[model.CompiledPolicy] {
				return repository.New[model.CompiledPolicy](db)
			},

			func(db *gorm.DB) repository.Base[model.Policy] {
				return repository.New[model.Policy](db)
			},

			func(db *gorm.DB) repository.Base[model.Principal] {
				return repository.New[model.Principal](db)
			},

			func(db *gorm.DB) repository.Base[model.Resource] {
				return repository.New[model.Resource](db)
			},

			func(db *gorm.DB) repository.Base[model.Role] {
				return repository.New[model.Role](db)
			},

			func(db *gorm.DB) repository.Base[model.Stats] {
				return repository.New[model.Stats](db)
			},

			func(db *gorm.DB) repository.Base[model.User] {
				return repository.New[model.User](db)
			},

			func(base repository.Base[model.Principal]) repository.Principal {
				return repository.NewPrincipal(base)
			},

			func(base repository.Base[model.Resource]) repository.Resource {
				return repository.NewResource(base)
			},
		),
	)
}
