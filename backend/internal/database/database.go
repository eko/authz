package database

import (
	"github.com/eko/authz/backend/configs"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *configs.Database, slogLogger *slog.Logger) (db *gorm.DB, err error) {
	logger := logger(slogLogger)
	logger.ignoreRecordNotFoundError = true

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DSN(),
	}), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	// TODO: delete these auto-migrate lines.
	// db.AutoMigrate(model.Action{})
	// db.AutoMigrate(model.Compiled{})
	// db.AutoMigrate(model.Resource{})
	// db.AutoMigrate(model.Policy{})
	// db.AutoMigrate(model.Role{})
	// db.AutoMigrate(model.Principal{})

	return db, nil
}
