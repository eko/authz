package interceptor

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/eko/authz/backend/pkg/authz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

type fakeServer struct {
	authz.UnimplementedApiServer

	authenticateCalls atomic.Int32

	mu             sync.Mutex
	authorizations [][]string
	forwarded      [][]string
}

func (s *fakeServer) Authenticate(context.Context, *authz.AuthenticateRequest) (*authz.AuthenticateResponse, error) {
	s.authenticateCalls.Add(1)

	return &authz.AuthenticateResponse{Token: "token-123", ExpiresIn: 3600}, nil
}

func (s *fakeServer) PrincipalGet(ctx context.Context, _ *authz.PrincipalGetRequest) (*authz.PrincipalGetResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.authorizations = append(s.authorizations, md.Get(authorizationKey))
	s.forwarded = append(s.forwarded, md.Get("x-incoming-only"))

	return &authz.PrincipalGetResponse{}, nil
}

func newTestClient(t *testing.T, server *fakeServer) authz.ApiClient {
	t.Helper()

	listener := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	authz.RegisterApiServer(grpcServer, server)

	go func() { _ = grpcServer.Serve(listener) }()
	t.Cleanup(grpcServer.Stop)

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return listener.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(AuthenticationUnaryClientInterceptor(NewAuthenticator("id", "secret"))),
	)
	if err != nil {
		t.Fatalf("cannot create client: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	return authz.NewApiClient(conn)
}

func TestAuthenticationUnaryClientInterceptor_ReusesToken(t *testing.T) {
	server := &fakeServer{}
	client := newTestClient(t, server)

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := client.PrincipalGet(context.Background(), &authz.PrincipalGetRequest{Id: "p"}); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()

	// Concurrent first calls may each authenticate, but later calls must reuse the token.
	before := server.authenticateCalls.Load()
	if _, err := client.PrincipalGet(context.Background(), &authz.PrincipalGetRequest{Id: "p"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if after := server.authenticateCalls.Load(); after != before {
		t.Errorf("expected cached token to be reused, Authenticate called %d more time(s)", after-before)
	}

	for _, authorization := range server.authorizations {
		if len(authorization) != 1 || authorization[0] != "bearer token-123" {
			t.Errorf("unexpected authorization metadata: %v", authorization)
		}
	}
}

func TestAuthenticationUnaryClientInterceptor_DoesNotForwardIncomingMetadata(t *testing.T) {
	server := &fakeServer{}
	client := newTestClient(t, server)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-incoming-only", "secret"))
	if _, err := client.PrincipalGet(ctx, &authz.PrincipalGetRequest{Id: "p"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if forwarded := server.forwarded[0]; len(forwarded) != 0 {
		t.Errorf("incoming metadata must not be forwarded, got %v", forwarded)
	}
}
