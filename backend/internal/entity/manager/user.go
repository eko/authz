package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/helper/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User interface {
	Create(username string, password string) (*model.User, error)
	Delete(username string) error
	GetRepository() repository.Base[model.User]
}

type userManager struct {
	repository          repository.Base[model.User]
	principalRepository repository.Base[model.Principal]
	transactionManager  database.TransactionManager
	tokenGenerator      token.Generator
}

// NewUser initializes a new user manager.
func NewUser(
	repository repository.Base[model.User],
	principalRepository repository.Base[model.Principal],
	transactionManager database.TransactionManager,
	tokenGenerator token.Generator,
) User {
	return &userManager{
		repository:          repository,
		principalRepository: principalRepository,
		transactionManager:  transactionManager,
		tokenGenerator:      tokenGenerator,
	}
}

func (m *userManager) GetRepository() repository.Base[model.User] {
	return m.repository
}

func (m *userManager) Create(username string, password string) (*model.User, error) {
	exists, err := m.repository.GetByFields(map[string]repository.FieldValue{
		"username": {Operator: "=", Value: username},
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing user: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a user already exists with username %q", username)
	}

	if password == "" {
		password, err = m.tokenGenerator.Generate(10)
		if err != nil {
			return nil, fmt.Errorf("unable to generate a random password: %v", err)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	transaction := m.transactionManager.New()
	defer func() { _ = transaction.Commit() }()

	user := &model.User{
		Username:     username,
		Password:     password,
		PasswordHash: string(hashedPassword),
	}

	if err := m.repository.WithTransaction(transaction).Create(user); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to create user: %v", err)
	}

	if err := m.principalRepository.WithTransaction(transaction).Create(&model.Principal{
		ID: model.UserPrincipal(user.Username),
	}); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to create user: %v", err)
	}

	return user, nil
}

func (m *userManager) Delete(username string) error {
	user, err := m.GetRepository().GetByFields(map[string]repository.FieldValue{
		"username": {Operator: "=", Value: username},
	})
	if err != nil {
		return fmt.Errorf("cannot retrieve user: %v", err)
	}

	// Retrieve principal
	principal, err := m.principalRepository.Get(
		model.UserPrincipal(username),
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve user principal: %v", err)
	}

	// Delete both user and principal
	transaction := m.transactionManager.New()
	defer func() { _ = transaction.Commit() }()

	if err := m.principalRepository.WithTransaction(transaction).Delete(principal); err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("cannot delete principal: %v", err)
	}

	if err := m.GetRepository().WithTransaction(transaction).Delete(user); err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("cannot delete user: %v", err)
	}

	return nil
}
