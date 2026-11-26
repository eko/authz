package interceptor

import (
	"sync"
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

// IsValid returns true when the token is still valid at the given time.
func (t *Token) IsValid(at time.Time) bool {
	return t != nil && t.AccessToken != "" && t.ExpireAt.After(at)
}

type Authenticator interface {
	GetClientID() string
	GetClientSecret() string
	GetExpireDelay() time.Duration
	GetToken() *Token
	SetToken(token *Token)
}

type authenticator struct {
	mu           sync.RWMutex
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
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.token
}

func (a *authenticator) SetToken(token *Token) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.token = token
}
