package handler

import (
	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/manager"
	"github.com/eko/authz/backend/internal/security/paseto"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

const (
	ActionGetKey        = "action-get"
	ActionListKey       = "action-list"
	AuthAuthenticateKey = "auth-authenticate"
	AuthTokenNewKey     = "auth-token-new"
	CheckKey            = "check"
	ClientCreateKey     = "client-create-key"
	ClientDeleteKey     = "client-delete-key"
	ClientGetKey        = "client-get-key"
	ClientListKey       = "client-list-key"
	PolicyCreateKey     = "policy-create"
	PolicyDeleteKey     = "policy-delete"
	PolicyGetKey        = "policy-get"
	PolicyListKey       = "policy-list"
	PolicyUpdateKey     = "policy-update"
	PrincipalCreateKey  = "principal-create"
	PrincipalDeleteKey  = "principal-delete"
	PrincipalGetKey     = "principal-get"
	PrincipalListKey    = "principal-list"
	PrincipalUpdateKey  = "principal-update"
	ResourceCreateKey   = "resource-create"
	ResourceDeleteKey   = "resource-delete"
	ResourceGetKey      = "resource-get"
	ResourceListKey     = "resource-list"
	ResourceUpdateKey   = "resource-update"
	RoleCreateKey       = "role-create"
	RoleDeleteKey       = "role-delete"
	RoleGetKey          = "role-get"
	RoleListKey         = "role-list"
	RoleUpdateKey       = "role-update"
	UserCreateKey       = "user-create-key"
	UserDeleteKey       = "user-delete-key"
	UserGetKey          = "user-get-key"
	UserListKey         = "user-list-key"
)

type Handler fiber.Handler
type Handlers map[string]Handler

func (h Handlers) Get(name string) Handler {
	return h[name]
}

func NewHandlers(
	authCfg *configs.Auth,
	validate *validator.Validate,
	manager manager.Manager,
	transactionManager database.TransactionManager,
	tokenManager paseto.Manager,
	oauthServer *server.Server,
) Handlers {
	return Handlers{
		ActionGetKey:        ActionGet(manager),
		ActionListKey:       ActionList(manager),
		AuthAuthenticateKey: Authenticate(validate, manager, authCfg, tokenManager),
		AuthTokenNewKey:     adaptor.HTTPHandlerFunc(TokenNew(oauthServer)),
		CheckKey:            Check(validate, manager),
		ClientCreateKey:     ClientCreate(validate, manager, authCfg),
		ClientDeleteKey:     ClientDelete(manager, transactionManager),
		ClientGetKey:        ClientGet(manager),
		ClientListKey:       ClientList(manager),
		PolicyCreateKey:     PolicyCreate(validate, manager),
		PolicyDeleteKey:     PolicyDelete(manager),
		PolicyGetKey:        PolicyGet(manager),
		PolicyListKey:       PolicyList(manager),
		PolicyUpdateKey:     PolicyUpdate(validate, manager),
		PrincipalCreateKey:  PrincipalCreate(validate, manager),
		PrincipalDeleteKey:  PrincipalDelete(manager),
		PrincipalGetKey:     PrincipalGet(manager),
		PrincipalListKey:    PrincipalList(manager),
		PrincipalUpdateKey:  PrincipalUpdate(validate, manager),
		ResourceCreateKey:   ResourceCreate(validate, manager),
		ResourceDeleteKey:   ResourceDelete(manager),
		ResourceGetKey:      ResourceGet(manager),
		ResourceListKey:     ResourceList(manager),
		ResourceUpdateKey:   ResourceUpdate(validate, manager),
		RoleCreateKey:       RoleCreate(validate, manager),
		RoleDeleteKey:       RoleDelete(manager),
		RoleGetKey:          RoleGet(manager),
		RoleListKey:         RoleList(manager),
		RoleUpdateKey:       RoleUpdate(validate, manager),
		UserCreateKey:       UserCreate(validate, manager),
		UserDeleteKey:       UserDelete(manager, transactionManager),
		UserGetKey:          UserGet(manager),
		UserListKey:         UserList(manager),
	}
}
