package paseto

import (
	lib_time "time"

	"aidanwoods.dev/go-paseto"
	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/helper/time"
)

const (
	TokenTypeBearer = "bearer"
)

type Token struct {
	Token     string
	TokenType string
	ExpiresIn int64
}

type Manager interface {
	Generate(identifier string, expiration lib_time.Duration) *Token
	Parse(token string) (map[string]any, error)
}

type manager struct {
	cfg       *configs.Auth
	clock     time.Clock
	secretKey paseto.V4AsymmetricSecretKey
}

func NewManager(
	cfg *configs.Auth,
	clock time.Clock,
) (Manager, error) {
	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(cfg.SecretKeyHex)
	if err != nil {
		return nil, err
	}

	return &manager{
		cfg:       cfg,
		clock:     clock,
		secretKey: secretKey,
	}, nil
}

func (g *manager) Generate(identifier string, expiration lib_time.Duration) *Token {
	token := paseto.NewToken()

	now := g.clock.Now()
	expireAt := now.Add(expiration)

	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(expireAt)
	token.SetString("id", identifier)

	return &Token{
		Token:     token.V4Sign(g.secretKey, nil),
		TokenType: TokenTypeBearer,
		ExpiresIn: int64(expireAt.Sub(g.clock.Now()).Seconds()),
	}
}

func (g *manager) Parse(token string) (map[string]any, error) {
	parser := paseto.NewParser()

	pasetoToken, err := parser.ParseV4Public(g.secretKey.Public(), token, nil)
	if err != nil {
		return nil, err
	}

	return pasetoToken.Claims(), nil
}
