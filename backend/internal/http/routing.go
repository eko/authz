package http

import (
	"strings"

	_ "github.com/eko/authz/backend/internal/http/docs"
	"github.com/eko/authz/backend/internal/http/handler"
	"github.com/eko/authz/backend/internal/http/middleware"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Is to define the swagger route and the dynamic swagger routes
func (s *Server) setSwagger() {
	s.app.Get("/swagger/*", swagger.HandlerDefault)
}

// @title						Authz API
// @version					1.0
// @description				Authorization management HTTP APIs
// @securitydefinitions.apikey	Authentication
// @in							header
// @name						Authorization
// @BasePath					/v1
func (s *Server) setRoutes() {
	s.app.Use(
		cors.New(cors.Config{
			AllowOrigins:     strings.Join(s.cfg.CORSAllowedDomains, ","),
			AllowMethods:     strings.Join(s.cfg.CORSAllowedMethods, ","),
			AllowHeaders:     strings.Join(s.cfg.CORSAllowedHeaders, ","),
			AllowCredentials: s.cfg.CORSAllowCredentials,
			MaxAge:           int(s.cfg.CORSCacheMaxAge.Seconds()),
		}),
		otelfiber.Middleware(),
	)

	base := s.app.Group("/v1")
	{
		// Authentication
		base.Post("/auth", s.handlers.Get(handler.AuthAuthenticateKey))
		base.Get("/oauth", s.handlers.Get(handler.OAuthAuthenticateKey))
		base.Get("/oauth/callback", s.handlers.Get(handler.OAuthCallbackKey))
		base.Post("/token", s.handlers.Get(handler.AuthTokenNewKey))

		// Observability
		base.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

		// Authz resources
		authenticated := base.Use(s.middlewares.Get(middleware.AuthenticationKey))

		authenticated.Post("/check", s.handlers.Get(handler.CheckKey))

		actions := authenticated.Group("/actions")
		actions.Get("", s.authorized("authz.actions", "list", s.handlers.Get(handler.ActionListKey))...)
		actions.Get("/:identifier", s.authorized("authz.actions", "get", s.handlers.Get(handler.ActionGetKey))...)

		audits := authenticated.Group("/audits")
		audits.Get("", s.authorized("authz.audits", "get", s.handlers.Get(handler.AuditGetKey))...)

		clients := authenticated.Group("/clients")
		clients.Post("", s.authorized("authz.clients", "create", s.handlers.Get(handler.ClientCreateKey))...)
		clients.Get("", s.authorized("authz.clients", "list", s.handlers.Get(handler.ClientListKey))...)
		clients.Get("/:identifier", s.authorized("authz.clients", "get", s.handlers.Get(handler.ClientGetKey))...)
		clients.Delete("/:identifier", s.authorized("authz.clients", "delete", s.handlers.Get(handler.ClientDeleteKey))...)

		compiled := authenticated.Group("/compiled")
		compiled.Get("", s.authorized("authz.compiled", "list", s.handlers.Get(handler.CompiledListKey))...)

		policies := authenticated.Group("/policies")
		policies.Post("", s.authorized("authz.policies", "create", s.handlers.Get(handler.PolicyCreateKey))...)
		policies.Get("", s.authorized("authz.policies", "list", s.handlers.Get(handler.PolicyListKey))...)
		policies.Get("/:identifier", s.authorized("authz.policies", "get", s.handlers.Get(handler.PolicyGetKey))...)
		policies.Delete("/:identifier", s.authorized("authz.policies", "delete", s.handlers.Get(handler.PolicyDeleteKey))...)
		policies.Put("/:identifier", s.authorized("authz.policies", "update", s.handlers.Get(handler.PolicyUpdateKey))...)

		principals := authenticated.Group("/principals")
		principals.Post("", s.authorized("authz.principals", "create", s.handlers.Get(handler.PrincipalCreateKey))...)
		principals.Get("", s.authorized("authz.principals", "list", s.handlers.Get(handler.PrincipalListKey))...)
		principals.Get("/:identifier", s.authorized("authz.principals", "get", s.handlers.Get(handler.PrincipalGetKey))...)
		principals.Delete("/:identifier", s.authorized("authz.principals", "delete", s.handlers.Get(handler.PrincipalDeleteKey))...)
		principals.Put("/:identifier", s.authorized("authz.principals", "update", s.handlers.Get(handler.PrincipalUpdateKey))...)

		resources := authenticated.Group("/resources")
		resources.Post("", s.authorized("authz.resources", "create", s.handlers.Get(handler.ResourceCreateKey))...)
		resources.Get("", s.authorized("authz.resources", "list", s.handlers.Get(handler.ResourceListKey))...)
		resources.Get("/:identifier", s.authorized("authz.resources", "get", s.handlers.Get(handler.ResourceGetKey))...)
		resources.Delete("/:identifier", s.authorized("authz.resources", "delete", s.handlers.Get(handler.ResourceDeleteKey))...)
		resources.Put("/:identifier", s.authorized("authz.resources", "update", s.handlers.Get(handler.ResourceUpdateKey))...)

		role := authenticated.Group("/roles")
		role.Post("", s.authorized("authz.roles", "create", s.handlers.Get(handler.RoleCreateKey))...)
		role.Get("", s.authorized("authz.roles", "list", s.handlers.Get(handler.RoleListKey))...)
		role.Get("/:identifier", s.authorized("authz.roles", "get", s.handlers.Get(handler.RoleGetKey))...)
		role.Delete("/:identifier", s.authorized("authz.roles", "delete", s.handlers.Get(handler.RoleDeleteKey))...)
		role.Put("/:identifier", s.authorized("authz.roles", "update", s.handlers.Get(handler.RoleUpdateKey))...)

		stats := authenticated.Group("/stats")
		stats.Get("", s.authorized("authz.stats", "get", s.handlers.Get(handler.StatsGetKey))...)

		users := authenticated.Group("/users")
		users.Post("", s.authorized("authz.users", "create", s.handlers.Get(handler.UserCreateKey))...)
		users.Get("", s.authorized("authz.users", "list", s.handlers.Get(handler.UserListKey))...)
		users.Get("/:identifier", s.authorized("authz.users", "get", s.handlers.Get(handler.UserGetKey))...)
		users.Delete("/:identifier", s.authorized("authz.users", "delete", s.handlers.Get(handler.UserDeleteKey))...)
	}
}

func (s *Server) authorized(resourceKind string, action string, handler fiber.Handler) []fiber.Handler {
	var handlers = []fiber.Handler{
		middleware.DiscoverResource(resourceKind, action),
		s.middlewares.Get(middleware.AuthorizationKey),
		handler,
	}

	return handlers
}
