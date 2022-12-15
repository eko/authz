package handler

import (
	"github.com/eko/authz/backend/internal/manager"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const (
	ActionListKey      = "action-list"
	ActionGetKey       = "action-get"
	PolicyCreateKey    = "policy-create"
	PolicyDeleteKey    = "policy-delete"
	PolicyUpdateKey    = "policy-update"
	PolicyListKey      = "policy-list"
	PolicyGetKey       = "policy-get"
	ResourceCreateKey  = "resource-create"
	ResourceDeleteKey  = "resource-delete"
	ResourceUpdateKey  = "resource-update"
	ResourceListKey    = "resource-list"
	ResourceGetKey     = "resource-get"
	PrincipalCreateKey = "principal-create"
	PrincipalDeleteKey = "principal-delete"
	PrincipalUpdateKey = "principal-update"
	PrincipalListKey   = "principal-list"
	PrincipalGetKey    = "principal-get"
)

type Handler fiber.Handler

type Handlers map[string]Handler

func (h Handlers) Get(name string) Handler {
	return h[name]
}

func NewHandlers(
	validate *validator.Validate,
	manager manager.Manager,
) Handlers {
	return Handlers{
		ActionListKey:      ActionList(manager),
		ActionGetKey:       ActionGet(manager),
		PolicyCreateKey:    PolicyCreate(validate, manager),
		PolicyDeleteKey:    PolicyDelete(manager),
		PolicyUpdateKey:    PolicyUpdate(validate, manager),
		PolicyListKey:      PolicyList(manager),
		PolicyGetKey:       PolicyGet(manager),
		ResourceCreateKey:  ResourceCreate(validate, manager),
		ResourceDeleteKey:  ResourceDelete(manager),
		ResourceUpdateKey:  ResourceUpdate(validate, manager),
		ResourceListKey:    ResourceList(manager),
		ResourceGetKey:     ResourceGet(manager),
		PrincipalCreateKey: PrincipalCreate(validate, manager),
		PrincipalDeleteKey: PrincipalDelete(manager),
		PrincipalUpdateKey: PrincipalUpdate(validate, manager),
		PrincipalListKey:   PrincipalList(manager),
		PrincipalGetKey:    PrincipalGet(manager),
	}
}
