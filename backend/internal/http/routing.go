package http

import (
	_ "github.com/eko/authz/backend/internal/http/docs"
	"github.com/eko/authz/backend/internal/http/handler"
	"github.com/eko/authz/backend/internal/http/middleware"
	fiberotel "github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Is to define the swagger route and the dynamic swagger routes
func (s *Server) setSwagger() {
	s.app.Get("/swagger/*", swaggo.HandlerDefault)
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
			AllowOrigins:     s.cfg.CORSAllowedDomains,
			AllowMethods:     s.cfg.CORSAllowedMethods,
			AllowHeaders:     s.cfg.CORSAllowedHeaders,
			AllowCredentials: s.cfg.CORSAllowCredentials,
			MaxAge:           int(s.cfg.CORSCacheMaxAge.Seconds()),
		}),
		fiberotel.New(),
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
		s.authorized(actions, fiber.MethodGet, "", "authz.actions", "list", s.handlers.Get(handler.ActionListKey))
		s.authorized(actions, fiber.MethodGet, "/:identifier", "authz.actions", "get", s.handlers.Get(handler.ActionGetKey))

		audits := authenticated.Group("/audits")
		s.authorized(audits, fiber.MethodGet, "", "authz.audits", "get", s.handlers.Get(handler.AuditGetKey))

		clients := authenticated.Group("/clients")
		s.authorized(clients, fiber.MethodPost, "", "authz.clients", "create", s.handlers.Get(handler.ClientCreateKey))
		s.authorized(clients, fiber.MethodGet, "", "authz.clients", "list", s.handlers.Get(handler.ClientListKey))
		s.authorized(clients, fiber.MethodGet, "/:identifier", "authz.clients", "get", s.handlers.Get(handler.ClientGetKey))
		s.authorized(clients, fiber.MethodDelete, "/:identifier", "authz.clients", "delete", s.handlers.Get(handler.ClientDeleteKey))

		compiled := authenticated.Group("/compiled")
		s.authorized(compiled, fiber.MethodGet, "", "authz.compiled", "list", s.handlers.Get(handler.CompiledListKey))

		policies := authenticated.Group("/policies")
		s.authorized(policies, fiber.MethodPost, "", "authz.policies", "create", s.handlers.Get(handler.PolicyCreateKey))
		s.authorized(policies, fiber.MethodGet, "", "authz.policies", "list", s.handlers.Get(handler.PolicyListKey))
		s.authorized(policies, fiber.MethodGet, "/:identifier", "authz.policies", "get", s.handlers.Get(handler.PolicyGetKey))
		s.authorized(policies, fiber.MethodDelete, "/:identifier", "authz.policies", "delete", s.handlers.Get(handler.PolicyDeleteKey))
		s.authorized(policies, fiber.MethodPut, "/:identifier", "authz.policies", "update", s.handlers.Get(handler.PolicyUpdateKey))

		principals := authenticated.Group("/principals")
		s.authorized(principals, fiber.MethodPost, "", "authz.principals", "create", s.handlers.Get(handler.PrincipalCreateKey))
		s.authorized(principals, fiber.MethodGet, "", "authz.principals", "list", s.handlers.Get(handler.PrincipalListKey))
		s.authorized(principals, fiber.MethodGet, "/:identifier", "authz.principals", "get", s.handlers.Get(handler.PrincipalGetKey))
		s.authorized(principals, fiber.MethodDelete, "/:identifier", "authz.principals", "delete", s.handlers.Get(handler.PrincipalDeleteKey))
		s.authorized(principals, fiber.MethodPut, "/:identifier", "authz.principals", "update", s.handlers.Get(handler.PrincipalUpdateKey))

		resources := authenticated.Group("/resources")
		s.authorized(resources, fiber.MethodPost, "", "authz.resources", "create", s.handlers.Get(handler.ResourceCreateKey))
		s.authorized(resources, fiber.MethodGet, "", "authz.resources", "list", s.handlers.Get(handler.ResourceListKey))
		s.authorized(resources, fiber.MethodGet, "/:identifier", "authz.resources", "get", s.handlers.Get(handler.ResourceGetKey))
		s.authorized(resources, fiber.MethodDelete, "/:identifier", "authz.resources", "delete", s.handlers.Get(handler.ResourceDeleteKey))
		s.authorized(resources, fiber.MethodPut, "/:identifier", "authz.resources", "update", s.handlers.Get(handler.ResourceUpdateKey))

		role := authenticated.Group("/roles")
		s.authorized(role, fiber.MethodPost, "", "authz.roles", "create", s.handlers.Get(handler.RoleCreateKey))
		s.authorized(role, fiber.MethodGet, "", "authz.roles", "list", s.handlers.Get(handler.RoleListKey))
		s.authorized(role, fiber.MethodGet, "/:identifier", "authz.roles", "get", s.handlers.Get(handler.RoleGetKey))
		s.authorized(role, fiber.MethodDelete, "/:identifier", "authz.roles", "delete", s.handlers.Get(handler.RoleDeleteKey))
		s.authorized(role, fiber.MethodPut, "/:identifier", "authz.roles", "update", s.handlers.Get(handler.RoleUpdateKey))

		stats := authenticated.Group("/stats")
		s.authorized(stats, fiber.MethodGet, "", "authz.stats", "get", s.handlers.Get(handler.StatsGetKey))

		users := authenticated.Group("/users")
		s.authorized(users, fiber.MethodPost, "", "authz.users", "create", s.handlers.Get(handler.UserCreateKey))
		s.authorized(users, fiber.MethodGet, "", "authz.users", "list", s.handlers.Get(handler.UserListKey))
		s.authorized(users, fiber.MethodGet, "/:identifier", "authz.users", "get", s.handlers.Get(handler.UserGetKey))
		s.authorized(users, fiber.MethodDelete, "/:identifier", "authz.users", "delete", s.handlers.Get(handler.UserDeleteKey))
	}
}

// authorized registers a route that requires the authenticated user
// to be allowed to perform the given action on the given resource kind.
func (s *Server) authorized(
	router fiber.Router,
	method string,
	path string,
	resourceKind string,
	action string,
	handler fiber.Handler,
) {
	router.Add(
		[]string{method},
		path,
		middleware.DiscoverResource(resourceKind, action),
		s.middlewares.Get(middleware.AuthorizationKey),
		handler,
	)
}
