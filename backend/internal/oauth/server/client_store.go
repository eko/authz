package server

import (
	"context"
	"encoding/json"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

type ClientStore struct {
	clientManager manager.Client
}

func NewClientStore(clientManager manager.Client) *ClientStore {
	return &ClientStore{
		clientManager: clientManager,
	}
}

func (s *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	if id == "" {
		return nil, nil
	}

	client, err := s.clientManager.GetRepository().Get(id)
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

	return s.clientManager.GetRepository().Create(&model.Client{
		ID:     info.GetID(),
		Secret: info.GetSecret(),
		Domain: info.GetDomain(),
		Data:   string(data),
	})
}
