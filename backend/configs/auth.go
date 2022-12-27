package configs

import "time"

type Auth struct {
	AccessTokenDuration  time.Duration `config:"auth_access_token_duration"`
	RefreshTokenDuration time.Duration `config:"auth_refresh_token_duration"`
	Domain               string        `config:"auth_domain"`
	SecretKeyHex         string        `config:"auth_secret_key_hex"`
}

func newAuth() *Auth {
	return &Auth{
		AccessTokenDuration:  6 * time.Hour,
		RefreshTokenDuration: 2 * time.Hour,
		Domain:               "http://localhost:8080",
		SecretKeyHex:         "7b65f3785060d7dcba351c4ddf1e864f0becd7e989f096c4a6ea28132dbb18cc08474e6caa2da411cff0be55f44f73714282e90bfd77e9c3859e0284f4098df8",
	}
}
