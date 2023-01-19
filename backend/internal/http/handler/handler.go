package handler

import (
	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/helper/token"
	"github.com/eko/authz/backend/internal/oauth/client"
	"github.com/eko/authz/backend/internal/security/jwt"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

const (
	ActionGetKey         = "action-get"
	ActionListKey        = "action-list"
	AuditGetKey          = "audit-get"
	AuthAuthenticateKey  = "auth-authenticate"
	AuthTokenNewKey      = "auth-token-new"
	CheckKey             = "check"
	ClientCreateKey      = "client-create"
	ClientDeleteKey      = "client-delete"
	ClientGetKey         = "client-get"
	ClientListKey        = "client-list"
	CompiledListKey      = "compiled-list"
	OAuthAuthenticateKey = "oauth-authenticate"
	OAuthCallbackKey     = "oauth-callback"
	PolicyCreateKey      = "policy-create"
	PolicyDeleteKey      = "policy-delete"
	PolicyGetKey         = "policy-get"
	PolicyListKey        = "policy-list"
	PolicyUpdateKey      = "policy-update"
	PrincipalCreateKey   = "principal-create"
	PrincipalDeleteKey   = "principal-delete"
	PrincipalGetKey      = "principal-get"
	PrincipalListKey     = "principal-list"
	PrincipalUpdateKey   = "principal-update"
	ResourceCreateKey    = "resource-create"
	ResourceDeleteKey    = "resource-delete"
	ResourceGetKey       = "resource-get"
	ResourceListKey      = "resource-list"
	ResourceUpdateKey    = "resource-update"
	RoleCreateKey        = "role-create"
	RoleDeleteKey        = "role-delete"
	RoleGetKey           = "role-get"
	RoleListKey          = "role-list"
	RoleUpdateKey        = "role-update"
	StatsGetKey          = "stats-get"
	UserCreateKey        = "user-create"
	UserDeleteKey        = "user-delete"
	UserGetKey           = "user-get"
	UserListKey          = "user-list"
)

type Handler fiber.Handler
type Handlers map[string]Handler

func (h Handlers) Get(name string) Handler {
	return h[name]
}

func NewHandlers(
	actionManager manager.Action,
	auditManager manager.Audit,
	authCfg *configs.Auth,
	clientManager manager.Client,
	compiledManager manager.CompiledPolicy,
	dispatcher event.Dispatcher,
	logger *slog.Logger,
	oauthClientManager client.Manager,
	oauthServer *server.Server,
	policyManager manager.Policy,
	principalManager manager.Principal,
	resourceManager manager.Resource,
	roleManager manager.Role,
	statsManager manager.Stats,
	tokenGenerator token.Generator,
	jwtManager jwt.Manager,
	userManager manager.User,
	validate *validator.Validate,
) Handlers {
	return Handlers{
		ActionGetKey:         ActionGet(actionManager),
		ActionListKey:        ActionList(actionManager),
		AuditGetKey:          AuditGet(auditManager),
		AuthAuthenticateKey:  Authenticate(validate, userManager, jwtManager),
		AuthTokenNewKey:      adaptor.HTTPHandlerFunc(TokenNew(oauthServer)),
		CheckKey:             Check(logger, validate, compiledManager, dispatcher),
		ClientCreateKey:      ClientCreate(validate, clientManager, authCfg),
		ClientDeleteKey:      ClientDelete(clientManager),
		ClientGetKey:         ClientGet(clientManager),
		ClientListKey:        ClientList(clientManager),
		CompiledListKey:      CompiledList(compiledManager),
		OAuthAuthenticateKey: OAuthAuthenticate(oauthClientManager, tokenGenerator),
		OAuthCallbackKey:     OAuthCallback(jwtManager, oauthClientManager, principalManager),
		PolicyCreateKey:      PolicyCreate(validate, policyManager),
		PolicyDeleteKey:      PolicyDelete(policyManager),
		PolicyGetKey:         PolicyGet(policyManager),
		PolicyListKey:        PolicyList(policyManager),
		PolicyUpdateKey:      PolicyUpdate(validate, policyManager),
		PrincipalCreateKey:   PrincipalCreate(validate, principalManager),
		PrincipalDeleteKey:   PrincipalDelete(principalManager),
		PrincipalGetKey:      PrincipalGet(principalManager),
		PrincipalListKey:     PrincipalList(principalManager),
		PrincipalUpdateKey:   PrincipalUpdate(validate, principalManager),
		ResourceCreateKey:    ResourceCreate(validate, resourceManager),
		ResourceDeleteKey:    ResourceDelete(resourceManager),
		ResourceGetKey:       ResourceGet(resourceManager),
		ResourceListKey:      ResourceList(resourceManager),
		ResourceUpdateKey:    ResourceUpdate(validate, resourceManager),
		RoleCreateKey:        RoleCreate(validate, roleManager),
		RoleDeleteKey:        RoleDelete(roleManager),
		RoleGetKey:           RoleGet(roleManager),
		RoleListKey:          RoleList(roleManager),
		RoleUpdateKey:        RoleUpdate(validate, roleManager),
		StatsGetKey:          StatsGet(statsManager),
		UserCreateKey:        UserCreate(validate, userManager),
		UserDeleteKey:        UserDelete(userManager),
		UserGetKey:           UserGet(userManager),
		UserListKey:          UserList(userManager),
	}
}
