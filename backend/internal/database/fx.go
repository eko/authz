package database

import (
	"github.com/eko/authz/backend/internal/database/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func FxModule() fx.Option {
	return fx.Module("database",
		fx.Provide(
			New,

			func(db *gorm.DB) *Repository[model.Action] {
				return NewRepository[model.Action](db)
			},

			func(db *gorm.DB) *Repository[model.CompiledPolicy] {
				return NewRepository[model.CompiledPolicy](db)
			},

			func(db *gorm.DB) *Repository[model.Policy] {
				return NewRepository[model.Policy](db)
			},

			func(db *gorm.DB) *Repository[model.Resource] {
				return NewRepository[model.Resource](db)
			},

			func(db *gorm.DB) *Repository[model.Role] {
				return NewRepository[model.Role](db)
			},

			func(db *gorm.DB) *Repository[model.Principal] {
				return NewRepository[model.Principal](db)
			},
		),
	)
}
