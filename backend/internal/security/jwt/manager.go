package jwt

import (
	"fmt"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/helper/time"
	"github.com/golang-jwt/jwt/v4"
)

const (
	TokenTypeBearer = "bearer"
)

type Claims struct {
	jwt.RegisteredClaims
}

type Token struct {
	Token     string
	TokenType string
	ExpiresIn int64
}

type Manager interface {
	Generate(identifier string) (*Token, error)
	Parse(accessToken string) (*Claims, error)
}

type manager struct {
	cfg   *configs.Auth
	clock time.Clock
}

func NewManager(
	cfg *configs.Auth,
	clock time.Clock,
) Manager {
	// Ensure JWT library is using our clock.
	jwt.TimeFunc = clock.Now

	return &manager{
		cfg:   cfg,
		clock: clock,
	}
}

func (g *manager) Generate(identifier string) (*Token, error) {
	now := g.clock.Now()
	expireAt := now.Add(g.cfg.AccessTokenDuration)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    configs.ApplicationName,
			Subject:   identifier,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(g.cfg.JWTSignString)
	if err != nil {
		return nil, err
	}

	return &Token{
		Token:     accessToken,
		TokenType: TokenTypeBearer,
		ExpiresIn: int64(expireAt.Sub(g.clock.Now()).Seconds()),
	}, nil
}

func (g *manager) Parse(accessToken string) (*Claims, error) {
	claims := &Claims{}

	jwtToken, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return g.cfg.JWTSignString, nil
	})
	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return claims, nil
}
