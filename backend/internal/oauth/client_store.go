package oauth

import (
	"context"
	"encoding/json"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/database/model"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

type ClientStore struct {
	clientRepository *database.Repository[model.Client]
}

func NewClientStore(clientRepository *database.Repository[model.Client]) *ClientStore {
	return &ClientStore{
		clientRepository: clientRepository,
	}
}

func (s *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	if id == "" {
		return nil, nil
	}

	client, err := s.clientRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return &models.Client{
		ID:     client.ID,
		Secret: client.Secret,
		Domain: client.Domain,
	}, nil
}

func (s *ClientStore) Create(ctx context.Context, info oauth2.ClientInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return s.clientRepository.Create(&model.Client{
		ID:     info.GetID(),
		Secret: info.GetSecret(),
		Domain: info.GetDomain(),
		Data:   string(data),
	})
}
