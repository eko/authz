package grpc

import (
	"context"

	"github.com/eko/authz/backend/pkg/authz"
)

func (s *Server) Authenticate(ctx context.Context, req *authz.AuthenticateRequest) (*authz.AuthenticateResponse, error) {
	return s.authHandler.Authenticate(ctx, req)
}

func (s *Server) Check(ctx context.Context, req *authz.CheckRequest) (*authz.CheckResponse, error) {
	return s.checkHandler.Check(ctx, req)
}

func (s *Server) PolicyCreate(ctx context.Context, req *authz.PolicyCreateRequest) (*authz.PolicyCreateResponse, error) {
	return s.policyHandler.PolicyCreate(ctx, req)
}

func (s *Server) PolicyDelete(ctx context.Context, req *authz.PolicyDeleteRequest) (*authz.PolicyDeleteResponse, error) {
	return s.policyHandler.PolicyDelete(ctx, req)
}

func (s *Server) PolicyGet(ctx context.Context, req *authz.PolicyGetRequest) (*authz.PolicyGetResponse, error) {
	return s.policyHandler.PolicyGet(ctx, req)
}

func (s *Server) PolicyUpdate(ctx context.Context, req *authz.PolicyUpdateRequest) (*authz.PolicyUpdateResponse, error) {
	return s.policyHandler.PolicyUpdate(ctx, req)
}

func (s *Server) PrincipalCreate(ctx context.Context, req *authz.PrincipalCreateRequest) (*authz.PrincipalCreateResponse, error) {
	return s.principalHandler.PrincipalCreate(ctx, req)
}

func (s *Server) PrincipalDelete(ctx context.Context, req *authz.PrincipalDeleteRequest) (*authz.PrincipalDeleteResponse, error) {
	return s.principalHandler.PrincipalDelete(ctx, req)
}

func (s *Server) PrincipalGet(ctx context.Context, req *authz.PrincipalGetRequest) (*authz.PrincipalGetResponse, error) {
	return s.principalHandler.PrincipalGet(ctx, req)
}

func (s *Server) PrincipalUpdate(ctx context.Context, req *authz.PrincipalUpdateRequest) (*authz.PrincipalUpdateResponse, error) {
	return s.principalHandler.PrincipalUpdate(ctx, req)
}

func (s *Server) ResourceCreate(ctx context.Context, req *authz.ResourceCreateRequest) (*authz.ResourceCreateResponse, error) {
	return s.resourceHandler.ResourceCreate(ctx, req)
}

func (s *Server) ResourceDelete(ctx context.Context, req *authz.ResourceDeleteRequest) (*authz.ResourceDeleteResponse, error) {
	return s.resourceHandler.ResourceDelete(ctx, req)
}

func (s *Server) ResourceGet(ctx context.Context, req *authz.ResourceGetRequest) (*authz.ResourceGetResponse, error) {
	return s.resourceHandler.ResourceGet(ctx, req)
}

func (s *Server) ResourceUpdate(ctx context.Context, req *authz.ResourceUpdateRequest) (*authz.ResourceUpdateResponse, error) {
	return s.resourceHandler.ResourceUpdate(ctx, req)
}

func (s *Server) RoleCreate(ctx context.Context, req *authz.RoleCreateRequest) (*authz.RoleCreateResponse, error) {
	return s.roleHandler.RoleCreate(ctx, req)
}

func (s *Server) RoleDelete(ctx context.Context, req *authz.RoleDeleteRequest) (*authz.RoleDeleteResponse, error) {
	return s.roleHandler.RoleDelete(ctx, req)
}

func (s *Server) RoleGet(ctx context.Context, req *authz.RoleGetRequest) (*authz.RoleGetResponse, error) {
	return s.roleHandler.RoleGet(ctx, req)
}

func (s *Server) RoleUpdate(ctx context.Context, req *authz.RoleUpdateRequest) (*authz.RoleUpdateResponse, error) {
	return s.roleHandler.RoleUpdate(ctx, req)
}
