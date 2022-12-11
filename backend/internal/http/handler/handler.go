package handler

import (
	"github.com/eko/authz/backend/internal/manager"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const (
	ActionCreateKey   = "action-create"
	ActionDeleteKey   = "action-delete"
	ActionUpdateKey   = "action-update"
	ActionListKey     = "action-list"
	ActionGetKey      = "action-get"
	PolicyCreateKey   = "policy-create"
	PolicyDeleteKey   = "policy-delete"
	PolicyUpdateKey   = "policy-update"
	PolicyListKey     = "policy-list"
	PolicyGetKey      = "policy-get"
	ResourceCreateKey = "resource-create"
	ResourceDeleteKey = "resource-delete"
	ResourceUpdateKey = "resource-update"
	ResourceListKey   = "resource-list"
	ResourceGetKey    = "resource-get"
	SubjectCreateKey  = "subject-create"
	SubjectDeleteKey  = "subject-delete"
	SubjectUpdateKey  = "subject-update"
	SubjectListKey    = "subject-list"
	SubjectGetKey     = "subject-get"
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
		ActionCreateKey:   ActionCreate(validate, manager),
		ActionDeleteKey:   ActionDelete(manager),
		ActionUpdateKey:   ActionUpdate(validate, manager),
		ActionListKey:     ActionList(manager),
		ActionGetKey:      ActionGet(manager),
		PolicyCreateKey:   PolicyCreate(validate, manager),
		PolicyDeleteKey:   PolicyDelete(manager),
		PolicyUpdateKey:   PolicyUpdate(validate, manager),
		PolicyListKey:     PolicyList(manager),
		PolicyGetKey:      PolicyGet(manager),
		ResourceCreateKey: ResourceCreate(validate, manager),
		ResourceDeleteKey: ResourceDelete(manager),
		ResourceUpdateKey: ResourceUpdate(validate, manager),
		ResourceListKey:   ResourceList(manager),
		ResourceGetKey:    ResourceGet(manager),
		SubjectCreateKey:  SubjectCreate(validate, manager),
		SubjectDeleteKey:  SubjectDelete(manager),
		SubjectUpdateKey:  SubjectUpdate(validate, manager),
		SubjectListKey:    SubjectList(manager),
		SubjectGetKey:     SubjectGet(manager),
	}
}
