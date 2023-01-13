package configs

import "time"

type Auth struct {
	AccessTokenDuration  time.Duration `config:"auth_access_token_duration"`
	RefreshTokenDuration time.Duration `config:"auth_refresh_token_duration"`
	Domain               string        `config:"auth_domain"`
	JWTSignString        []byte        `config:"auth_jwt_sign_string"`
}

func newAuth() *Auth {
	return &Auth{
		AccessTokenDuration:  6 * time.Hour,
		RefreshTokenDuration: 6 * time.Hour,
		Domain:               "http://localhost:8080",
		JWTSignString:        []byte(`4uthz-s3cr3t-valu3-pl3as3-ch4ng3!`),
	}
}
