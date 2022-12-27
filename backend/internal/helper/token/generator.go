package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Generator is an interface that defines a token generator.
type Generator interface {
	Generate(length int) (string, error)
}

type generator struct{}

// NewGenerator returns a new token generator instance.
func NewGenerator() Generator {
	return &generator{}
}

// Generate generates a random string of length n.
func (*generator) Generate(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("token: can't generate a token of length %d", length)
	}

	token := make([]byte, length)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(token)[:length], nil
}
