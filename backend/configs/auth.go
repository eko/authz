package configs

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	// legacyJWTSignString is the signing secret that was shipped as default in
	// previous releases. It is publicly known and must never be accepted.
	legacyJWTSignString = `4uthz-s3cr3t-valu3-pl3as3-ch4ng3!`

	// minJWTSignStringLength is the minimum HMAC secret length (256 bits) for HS256.
	minJWTSignStringLength = 32

	// generatedJWTSignStringBytes random bytes give a 64 characters base64 secret.
	generatedJWTSignStringBytes = 48
)

var (
	ErrJWTSignStringLegacyDefault = errors.New("AUTH_JWT_SIGN_STRING must not use the publicly known legacy default value")
	ErrJWTSignStringTooShort      = fmt.Errorf("AUTH_JWT_SIGN_STRING must be at least %d bytes long", minJWTSignStringLength)
)

type Auth struct {
	AccessTokenDuration  time.Duration `config:"auth_access_token_duration"`
	RefreshTokenDuration time.Duration `config:"auth_refresh_token_duration"`
	Domain               string        `config:"auth_domain"`
	JWTSignString        string        `config:"auth_jwt_sign_string"`
}

func newAuth() *Auth {
	return &Auth{
		AccessTokenDuration:  6 * time.Hour,
		RefreshTokenDuration: 6 * time.Hour,
		Domain:               "http://localhost:8080",
	}
}

// ensureJWTSignString rejects weak signing secrets and generates a random one
// when none is provided, so tokens can never be signed with a known key.
func (a *Auth) ensureJWTSignString() error {
	if a.JWTSignString == "" {
		secret := make([]byte, generatedJWTSignStringBytes)
		if _, err := rand.Read(secret); err != nil {
			return fmt.Errorf("cannot generate JWT signing secret: %w", err)
		}

		a.JWTSignString = base64.RawURLEncoding.EncodeToString(secret)

		log.Print(
			"AUTH_JWT_SIGN_STRING is not set: a random signing secret has been generated. " +
				"Tokens will be invalidated on restart and will not be shared across instances. " +
				"Set AUTH_JWT_SIGN_STRING to a random value of at least 32 bytes in production.",
		)

		return nil
	}

	if a.JWTSignString == legacyJWTSignString {
		return ErrJWTSignStringLegacyDefault
	}

	if len(a.JWTSignString) < minJWTSignStringLength {
		return ErrJWTSignStringTooShort
	}

	return nil
}
