package manager

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"gorm.io/gorm"
)

type ActionRepository repository.Base[model.Action]

type Action interface {
	Create(identifier string) (*model.Action, error)
	GetRepository() ActionRepository
}

type actionManager struct {
	repository ActionRepository
}

// NewAction initializes a new action manager.
func NewAction(repository ActionRepository) Action {
	return &actionManager{
		repository: repository,
	}
}

func (m *actionManager) GetRepository() ActionRepository {
	return m.repository
}

func (m *actionManager) Create(identifier string) (*model.Action, error) {
	exists, err := m.repository.Get(identifier)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to check for existing action: %v", err)
	}

	if exists != nil {
		return nil, fmt.Errorf("an action already exists with identifier %q", identifier)
	}

	action := &model.Action{
		ID: identifier,
	}

	if err := m.repository.Create(action); err != nil {
		return nil, fmt.Errorf("unable to create action: %v", err)
	}

	return action, nil
}
