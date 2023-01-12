package handler

import (
	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/security/jwt"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
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
	StatsGetKey         = "stats-get"
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
	logger *slog.Logger,
	validate *validator.Validate,
	tokenManager jwt.Manager,
	dispatcher event.Dispatcher,
	actionManager manager.Action,
	clientManager manager.Client,
	compiledManager manager.CompiledPolicy,
	policyManager manager.Policy,
	principalManager manager.Principal,
	resourceManager manager.Resource,
	roleManager manager.Role,
	statsManager manager.Stats,
	userManager manager.User,
	oauthServer *server.Server,
) Handlers {
	return Handlers{
		ActionGetKey:        ActionGet(actionManager),
		ActionListKey:       ActionList(actionManager),
		AuthAuthenticateKey: Authenticate(validate, userManager, tokenManager),
		AuthTokenNewKey:     adaptor.HTTPHandlerFunc(TokenNew(oauthServer)),
		CheckKey:            Check(logger, validate, compiledManager, dispatcher),
		ClientCreateKey:     ClientCreate(validate, clientManager, authCfg),
		ClientDeleteKey:     ClientDelete(clientManager),
		ClientGetKey:        ClientGet(clientManager),
		ClientListKey:       ClientList(clientManager),
		PolicyCreateKey:     PolicyCreate(validate, policyManager),
		PolicyDeleteKey:     PolicyDelete(policyManager),
		PolicyGetKey:        PolicyGet(policyManager),
		PolicyListKey:       PolicyList(policyManager),
		PolicyUpdateKey:     PolicyUpdate(validate, policyManager),
		PrincipalCreateKey:  PrincipalCreate(validate, principalManager),
		PrincipalDeleteKey:  PrincipalDelete(principalManager),
		PrincipalGetKey:     PrincipalGet(principalManager),
		PrincipalListKey:    PrincipalList(principalManager),
		PrincipalUpdateKey:  PrincipalUpdate(validate, principalManager),
		ResourceCreateKey:   ResourceCreate(validate, resourceManager),
		ResourceDeleteKey:   ResourceDelete(resourceManager),
		ResourceGetKey:      ResourceGet(resourceManager),
		ResourceListKey:     ResourceList(resourceManager),
		ResourceUpdateKey:   ResourceUpdate(validate, resourceManager),
		RoleCreateKey:       RoleCreate(validate, roleManager),
		RoleDeleteKey:       RoleDelete(roleManager),
		RoleGetKey:          RoleGet(roleManager),
		RoleListKey:         RoleList(roleManager),
		RoleUpdateKey:       RoleUpdate(validate, roleManager),
		StatsGetKey:         StatsGet(statsManager),
		UserCreateKey:       UserCreate(validate, userManager),
		UserDeleteKey:       UserDelete(userManager),
		UserGetKey:          UserGet(userManager),
		UserListKey:         UserList(userManager),
	}
}
