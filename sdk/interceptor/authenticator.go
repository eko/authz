package interceptor

import (
	"time"
)

var (
	defaultExpireDelay = 5 * time.Minute
)

type AuthenticatorOption func(*authenticatorOptions)

type authenticatorOptions struct {
	expireDelay time.Duration
}

func WithExpireDelay(delay time.Duration) AuthenticatorOption {
	return func(o *authenticatorOptions) {
		o.expireDelay = delay
	}
}

type Token struct {
	AccessToken string
	ExpireAt    time.Time
}

type Authenticator interface {
	GetClientID() string
	GetClientSecret() string
	GetExpireDelay() time.Duration
	GetToken() *Token
	SetToken(token *Token)
}

type authenticator struct {
	clientID     string
	clientSecret string
	expireDelay  time.Duration
	token        *Token
}

func NewAuthenticator(clientID string, clientSecret string, options ...AuthenticatorOption) Authenticator {
	opts := &authenticatorOptions{
		expireDelay: defaultExpireDelay,
	}

	for _, o := range options {
		o(opts)
	}

	return &authenticator{
		clientID:     clientID,
		clientSecret: clientSecret,
		expireDelay:  opts.expireDelay,
	}
}

func (a *authenticator) GetClientID() string {
	return a.clientID
}

func (a *authenticator) GetClientSecret() string {
	return a.clientSecret
}

func (a *authenticator) GetExpireDelay() time.Duration {
	return a.expireDelay
}

func (a *authenticator) GetToken() *Token {
	return a.token
}

func (a *authenticator) SetToken(token *Token) {
	a.token = token
}
