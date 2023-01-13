package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/helper/token"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientRepository repository.Base[model.Client]

type Client interface {
	Create(name string, domain string) (*model.Client, error)
	Delete(identifier string) error
	GetRepository() ClientRepository
}

type clientManager struct {
	repository          ClientRepository
	principalRepository repository.Principal
	transactionManager  database.TransactionManager
	tokenGenerator      token.Generator
}

// NewClient initializes a new client manager.
func NewClient(
	repository ClientRepository,
	principalRepository repository.Principal,
	transactionManager database.TransactionManager,
	tokenGenerator token.Generator,
) Client {
	return &clientManager{
		repository:          repository,
		principalRepository: principalRepository,
		transactionManager:  transactionManager,
		tokenGenerator:      tokenGenerator,
	}
}

func (m *clientManager) GetRepository() ClientRepository {
	return m.repository
}

func (m *clientManager) Create(name string, domain string) (*model.Client, error) {
	exists, err := m.repository.GetByFields(map[string]repository.FieldValue{
		"name": {Operator: "=", Value: name},
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing client: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("a client already exists with name %q", name)
	}

	clientID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate client identifier: %v", err)
	}

	secret, err := m.tokenGenerator.Generate(48)
	if err != nil {
		return nil, fmt.Errorf("unable to generate client secret: %v", err)
	}

	client := &model.Client{
		ID:     clientID.String(),
		Secret: secret,
		Domain: domain,
		Name:   name,
	}

	transaction := m.transactionManager.New()
	defer func() { _ = transaction.Commit() }()

	if err := m.repository.WithTransaction(transaction).Create(client); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to create client: %v", err)
	}

	if err := m.principalRepository.WithTransaction(transaction).Create(&model.Principal{
		ID:       model.ClientPrincipal(client.Name),
		IsLocked: true,
	}); err != nil {
		_ = transaction.Rollback()
		return nil, fmt.Errorf("unable to create principal: %v", err)
	}

	return client, nil
}

func (m *clientManager) Delete(identifier string) error {
	client, err := m.GetRepository().Get(identifier)
	if err != nil {
		return fmt.Errorf("cannot retrieve client: %v", err)
	}

	// Retrieve principal
	principal, err := m.principalRepository.Get(
		model.ClientPrincipal(client.Name),
	)
	if err != nil {
		return fmt.Errorf("cannot retrieve client principal: %v", err)
	}

	// Delete both client and principal
	transaction := m.transactionManager.New()
	defer func() { _ = transaction.Commit() }()

	if err := m.principalRepository.WithTransaction(transaction).Delete(principal); err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("cannot delete principal: %v", err)
	}

	if err := m.GetRepository().WithTransaction(transaction).Delete(client); err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("cannot delete client: %v", err)
	}

	return nil
}
