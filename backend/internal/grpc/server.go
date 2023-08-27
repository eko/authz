package grpc

import (
	"context"
	"net"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/grpc/handler"
	"github.com/eko/authz/backend/internal/grpc/interceptor"
	"github.com/eko/authz/backend/internal/security/jwt"
	"github.com/eko/authz/backend/pkg/authz"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	authz.UnimplementedApiServer

	authHandler      handler.Auth
	checkHandler     handler.Check
	policyHandler    handler.Policy
	principalHandler handler.Principal
	resourceHandler  handler.Resource
	roleHandler      handler.Role

	addr       string
	GrpcServer *grpc.Server
}

func NewServer(
	cfg *configs.GRPCServer,
	tokenManager jwt.Manager,
	compiledManager manager.CompiledPolicy,
	authHandler handler.Auth,
	checkHandler handler.Check,
	policyHandler handler.Policy,
	principalHandler handler.Principal,
	resourceHandler handler.Resource,
	roleHandler handler.Role,
) *Server {
	server := &Server{
		addr:             cfg.Addr,
		authHandler:      authHandler,
		checkHandler:     checkHandler,
		policyHandler:    policyHandler,
		principalHandler: principalHandler,
		resourceHandler:  resourceHandler,
		roleHandler:      roleHandler,
	}

	authenticateFunc := interceptor.AuthenticateFunc(tokenManager)
	authorizationFunc := interceptor.AuthorizationFunc(compiledManager)

	grpcServer := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
			grpc_auth.StreamServerInterceptor(authenticateFunc),
			interceptor.AuthorizationStreamServerInterceptor(authorizationFunc),
		),
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			interceptor.AuthenticationUnaryServerInterceptor(
				grpc_auth.UnaryServerInterceptor(authenticateFunc),
			),
			interceptor.AuthorizationUnaryServerInterceptor(authorizationFunc),
		),
	)

	authz.RegisterApiServer(grpcServer, server)
	reflection.Register(grpcServer)

	server.GrpcServer = grpcServer

	return server
}

func Run(lc fx.Lifecycle, logger *slog.Logger, server *Server) error {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen("tcp", server.addr)
			if err != nil {
				return err
			}

			logger.Info("Starting gRPC server", slog.String("addr", server.addr))

			go func() {
				if err := server.GrpcServer.Serve(listener); err != nil {
					logger.Error("Unable to start gRPC server", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping gRPC server")

			server.GrpcServer.GracefulStop()
			return nil
		},
	})

	return nil
}
