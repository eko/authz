package database

import (
	"gorm.io/gorm"
)

type TransactionManager interface {
	New() Transaction
}

type Transaction interface {
	DB() *gorm.DB
	Commit() error
	Rollback() error
}

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{
		db: db,
	}
}

func (m *transactionManager) New() Transaction {
	return &transaction{
		db: m.db.Begin(),
	}
}

type transaction struct {
	db *gorm.DB
}

func (t *transaction) DB() *gorm.DB {
	return t.db
}

func (t *transaction) Commit() error {
	return t.db.Commit().Error
}

func (t *transaction) Rollback() error {
	return t.db.Rollback().Error
}
