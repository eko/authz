// Package fiber provides a GoFiber middleware that authorizes requests
// against Authz.
package fiber

import (
	"context"
	"errors"
	"fmt"

	"github.com/eko/authz/backend/pkg/authz"
	gofiber "github.com/gofiber/fiber/v3"
)

// WildcardValue is the resource value used when no specific resource is targeted.
const WildcardValue = "*"

var (
	// ErrUnauthenticated is passed to the error handler when the principal cannot be identified.
	ErrUnauthenticated = errors.New("unable to identify principal")

	// ErrForbidden is passed to the error handler when the principal is not allowed.
	ErrForbidden = errors.New("access denied")
)

// Checker performs authorization checks. The sdk.Client implements it.
type Checker interface {
	IsAllowed(ctx context.Context, check *authz.Check) (bool, error)
}

// Config configures the authorization middleware.
type Config struct {
	// Checker is used to perform the check. Required.
	Checker Checker

	// Principal returns the principal identifier of the current request. Required.
	Principal func(c gofiber.Ctx) (string, error)

	// ResourceKind is the kind of resource protected by the middleware. Required.
	ResourceKind string

	// ResourceValue returns the resource value targeted by the request.
	// Defaults to the wildcard value.
	ResourceValue func(c gofiber.Ctx) string

	// Action is the action performed on the resource. Required.
	Action string

	// ErrorHandler is called when the request is not authorized or when
	// an error occurs. Defaults to DefaultErrorHandler.
	ErrorHandler func(c gofiber.Ctx, err error) error
}

// New returns a middleware that only lets authorized requests through.
// It panics when the configuration is invalid.
func New(cfg Config) gofiber.Handler {
	if err := cfg.validate(); err != nil {
		panic(err)
	}

	if cfg.ResourceValue == nil {
		cfg.ResourceValue = func(gofiber.Ctx) string { return WildcardValue }
	}

	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = DefaultErrorHandler
	}

	return func(c gofiber.Ctx) error {
		principal, err := cfg.Principal(c)
		if err != nil || principal == "" {
			return cfg.ErrorHandler(c, fmt.Errorf("%w: %v", ErrUnauthenticated, err))
		}

		isAllowed, err := cfg.Checker.IsAllowed(c.Context(), &authz.Check{
			Principal:     principal,
			ResourceKind:  cfg.ResourceKind,
			ResourceValue: cfg.ResourceValue(c),
			Action:        cfg.Action,
		})
		if err != nil {
			return cfg.ErrorHandler(c, fmt.Errorf("unable to check authorization: %w", err))
		}

		if !isAllowed {
			return cfg.ErrorHandler(c, ErrForbidden)
		}

		return c.Next()
	}
}

// ResourceValueFromParam returns the value of the given route parameter,
// or the wildcard value when it is empty.
func ResourceValueFromParam(name string) func(c gofiber.Ctx) string {
	return func(c gofiber.Ctx) string {
		return c.Params(name, WildcardValue)
	}
}

// DefaultErrorHandler responds 401 when the principal is unknown, 403 when
// access is denied and 500 otherwise, without leaking error details.
func DefaultErrorHandler(c gofiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrUnauthenticated):
		return gofiber.ErrUnauthorized
	case errors.Is(err, ErrForbidden):
		return gofiber.ErrForbidden
	default:
		return gofiber.ErrInternalServerError
	}
}

func (cfg Config) validate() error {
	switch {
	case cfg.Checker == nil:
		return errors.New("authz fiber middleware: Checker is required")
	case cfg.Principal == nil:
		return errors.New("authz fiber middleware: Principal is required")
	case cfg.ResourceKind == "":
		return errors.New("authz fiber middleware: ResourceKind is required")
	case cfg.Action == "":
		return errors.New("authz fiber middleware: Action is required")
	default:
		return nil
	}
}
