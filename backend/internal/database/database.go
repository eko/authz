package database

import (
	"errors"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/helper/time"
	"golang.org/x/exp/slog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// ErrUnsupportedDriver is returned when specified database driver does not exists.
	ErrUnsupportedDriver = errors.New("unsupported database driver")
)

func New(
	cfg *configs.Database,
	slogLogger *slog.Logger,
	clock time.Clock,
) (db *gorm.DB, err error) {
	logger := logger(slogLogger)
	logger.ignoreRecordNotFoundError = true

	var dialector gorm.Dialector

	switch cfg.Driver {
	case configs.DriverMysql:
		dialector = mysql.Open(cfg.MysqlDSN())
	case configs.DriverSqlite:
		dialector = sqlite.Open(cfg.SqliteDSN())
	case configs.DriverPostgres:
		dialector = postgres.New(postgres.Config{DSN: cfg.PostgresDSN()})
	default:
		return nil, ErrUnsupportedDriver
	}

	db, err = gorm.Open(dialector, &gorm.Config{
		Logger:  logger,
		NowFunc: clock.Now,
	})
	if err != nil {
		return nil, err
	}

	if cfg.Driver == configs.DriverSqlite {
		checkErr(slogLogger, db.AutoMigrate(model.Action{}))
		checkErr(slogLogger, db.AutoMigrate(model.Attribute{}))
		checkErr(slogLogger, db.AutoMigrate(model.Client{}))
		checkErr(slogLogger, db.AutoMigrate(model.CompiledPolicy{}))
		checkErr(slogLogger, db.AutoMigrate(model.Policy{}))
		checkErr(slogLogger, db.AutoMigrate(model.Principal{}))
		checkErr(slogLogger, db.AutoMigrate(model.Stats{}))
		checkErr(slogLogger, db.AutoMigrate(model.Resource{}))
		checkErr(slogLogger, db.AutoMigrate(model.Role{}))
		checkErr(slogLogger, db.AutoMigrate(model.Token{}))
		checkErr(slogLogger, db.AutoMigrate(model.User{}))
	}

	return db, nil
}

func checkErr(logger *slog.Logger, err error) {
	if err != nil {
		logger.Error("Cannot migrate database", err)
	}
}
