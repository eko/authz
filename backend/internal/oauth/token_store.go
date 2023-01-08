package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"gorm.io/gorm"
)

const (
	gcInterval = 1000
)

func NewTokenStore(cfg *configs.Auth, db *gorm.DB) *TokenStore {
	store := &TokenStore{
		db:        db,
		tableName: model.Token{}.TableName(),
		stdout:    os.Stderr,
		ticker:    time.NewTicker(time.Second * time.Duration(gcInterval)),
	}

	go store.gc()

	return store
}

// Store mysql token store
type TokenStore struct {
	tableName string
	db        *gorm.DB
	stdout    io.Writer
	ticker    *time.Ticker
}

// SetStdout set error output
func (s *TokenStore) SetStdout(stdout io.Writer) *TokenStore {
	s.stdout = stdout
	return s
}

// Close close the store
func (s *TokenStore) Close() {
	s.ticker.Stop()
}

func (s *TokenStore) errorf(format string, args ...interface{}) {
	if s.stdout != nil {
		buf := fmt.Sprintf(format, args...)
		_, _ = s.stdout.Write([]byte(buf))
	}
}

func (s *TokenStore) gc() {
	for range s.ticker.C {
		now := time.Now().Unix()
		var count int64
		if err := s.db.Table(s.tableName).Where("expired_at <= ?", now).Or("code = ? and access = ? AND refresh = ?", "", "", "").Count(&count).Error; err != nil {
			s.errorf("[ERROR]:%s\n", err)
			return
		}
		if count > 0 {
			// not soft delete.
			if err := s.db.Table(s.tableName).Where("expired_at <= ?", now).Or("code = ? and access = ? AND refresh = ?", "", "", "").Unscoped().Delete(&model.Token{}).Error; err != nil {
				s.errorf("[ERROR]:%s\n", err)
			}
		}
	}
}

// Create create and store the new token information
func (s *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	jv, err := json.Marshal(info)
	if err != nil {
		return err
	}
	item := &model.Token{
		Data: string(jv),
	}

	if code := info.GetCode(); code != "" {
		item.Code = code
		item.ExpiredAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()).Unix()
	} else {
		item.Access = info.GetAccess()
		item.ExpiredAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).Unix()

		if refresh := info.GetRefresh(); refresh != "" {
			item.Refresh = info.GetRefresh()
			item.ExpiredAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Unix()
		}
	}

	return s.db.WithContext(ctx).Table(s.tableName).Create(item).Error
}

// RemoveByCode delete the authorization code
func (s *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	return s.db.WithContext(ctx).
		Table(s.tableName).
		Where("code = ?", code).
		Update("code", "").
		Error
}

// RemoveByAccess use the access token to delete the token information
func (s *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return s.db.WithContext(ctx).
		Table(s.tableName).
		Where("access = ?", access).
		Update("access", "").
		Error
}

// RemoveByRefresh use the refresh token to delete the token information
func (s *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return s.db.WithContext(ctx).
		Table(s.tableName).
		Where("refresh = ?", refresh).
		Update("refresh", "").
		Error
}

func (s *TokenStore) toTokenInfo(data string) oauth2.TokenInfo {
	var tm models.Token
	err := json.Unmarshal([]byte(data), &tm)
	if err != nil {
		return nil
	}
	return &tm
}

// GetByCode use the authorization code for token information data
func (s *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	if code == "" {
		return nil, nil
	}

	var item model.Token
	if err := s.db.WithContext(ctx).
		Table(s.tableName).
		Where("code = ?", code).
		Find(&item).Error; err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, nil
	}

	return s.toTokenInfo(item.Data), nil
}

// GetByAccess use the access token for token information data
func (s *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	if access == "" {
		return nil, nil
	}

	var item model.Token
	if err := s.db.WithContext(ctx).
		Table(s.tableName).
		Where("access = ?", access).
		Find(&item).Error; err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, nil
	}

	return s.toTokenInfo(item.Data), nil
}

// GetByRefresh use the refresh token for token information data
func (s *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	if refresh == "" {
		return nil, nil
	}

	var item model.Token
	if err := s.db.WithContext(ctx).
		Table(s.tableName).
		Where("refresh = ?", refresh).
		Find(&item).Error; err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, nil
	}

	return s.toTokenInfo(item.Data), nil
}
