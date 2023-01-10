package handler

import (
	"context"
	"fmt"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/entity/transformer"
	"github.com/eko/authz/backend/internal/http/handler/validator"
	"github.com/eko/authz/backend/pkg/authz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Role interface {
	RoleCreate(ctx context.Context, req *authz.RoleCreateRequest) (*authz.RoleCreateResponse, error)
	RoleDelete(ctx context.Context, req *authz.RoleDeleteRequest) (*authz.RoleDeleteResponse, error)
	RoleGet(ctx context.Context, req *authz.RoleGetRequest) (*authz.RoleGetResponse, error)
	RoleUpdate(ctx context.Context, req *authz.RoleUpdateRequest) (*authz.RoleUpdateResponse, error)
}

type role struct {
	roleManager manager.Role
}

func NewRole(
	roleManager manager.Role,
) Role {
	return &role{
		roleManager: roleManager,
	}
}

func (h *role) RoleCreate(ctx context.Context, req *authz.RoleCreateRequest) (*authz.RoleCreateResponse, error) {
	if !validator.ValidateSlugFromString(req.GetId()) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("identifier must be a slug, found: %s", req.GetId()))
	}

	role, err := h.roleManager.Create(req.GetId(), req.GetPolicies())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to create: %v", err.Error()))
	}

	return &authz.RoleCreateResponse{
		Role: transformer.NewRole(role).ToProto(),
	}, nil
}

func (h *role) RoleDelete(ctx context.Context, req *authz.RoleDeleteRequest) (*authz.RoleDeleteResponse, error) {
	err := h.roleManager.Delete(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to delete: %v", err.Error()))
	}

	return &authz.RoleDeleteResponse{
		Success: true,
	}, nil
}

func (h *role) RoleGet(ctx context.Context, req *authz.RoleGetRequest) (*authz.RoleGetResponse, error) {
	role, err := h.roleManager.GetRepository().Get(req.GetId(), repository.WithPreloads("Policies"))
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to retrieve: %v", err.Error()))
	}

	return &authz.RoleGetResponse{
		Role: transformer.NewRole(role).ToProto(),
	}, nil
}

func (h *role) RoleUpdate(ctx context.Context, req *authz.RoleUpdateRequest) (*authz.RoleUpdateResponse, error) {
	role, err := h.roleManager.Update(req.GetId(), req.GetPolicies())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to update: %v", err.Error()))
	}

	return &authz.RoleUpdateResponse{
		Role: transformer.NewRole(role).ToProto(),
	}, nil
}
