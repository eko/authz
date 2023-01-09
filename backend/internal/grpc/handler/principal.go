package handler

import (
	"context"
	"fmt"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/entity/transformer"
	"github.com/eko/authz/backend/pkg/authz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Principal interface {
	PrincipalCreate(ctx context.Context, req *authz.PrincipalCreateRequest) (*authz.PrincipalCreateResponse, error)
	PrincipalDelete(ctx context.Context, req *authz.PrincipalDeleteRequest) (*authz.PrincipalDeleteResponse, error)
	PrincipalGet(ctx context.Context, req *authz.PrincipalGetRequest) (*authz.PrincipalGetResponse, error)
	PrincipalUpdate(ctx context.Context, req *authz.PrincipalUpdateRequest) (*authz.PrincipalUpdateResponse, error)
}

type principal struct {
	principalManager manager.Principal
}

func NewPrincipal(
	principalManager manager.Principal,
) Principal {
	return &principal{
		principalManager: principalManager,
	}
}

func (h *principal) PrincipalCreate(ctx context.Context, req *authz.PrincipalCreateRequest) (*authz.PrincipalCreateResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	principal, err := h.principalManager.Create(req.GetId(), req.GetRoles(), attributesMap(req.GetAttributes()))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to create: %v", err.Error()))
	}

	return &authz.PrincipalCreateResponse{
		Principal: transformer.NewPrincipal(principal).ToProto(),
	}, nil
}

func (h *principal) PrincipalDelete(ctx context.Context, req *authz.PrincipalDeleteRequest) (*authz.PrincipalDeleteResponse, error) {
	err := h.principalManager.Delete(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to delete: %v", err.Error()))
	}

	return &authz.PrincipalDeleteResponse{
		Success: true,
	}, nil
}

func (h *principal) PrincipalGet(ctx context.Context, req *authz.PrincipalGetRequest) (*authz.PrincipalGetResponse, error) {
	principal, err := h.principalManager.GetRepository().Get(req.GetId(), repository.WithPreloads("Attributes", "Roles"))
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to retrieve: %v", err.Error()))
	}

	return &authz.PrincipalGetResponse{
		Principal: transformer.NewPrincipal(principal).ToProto(),
	}, nil
}

func (h *principal) PrincipalUpdate(ctx context.Context, req *authz.PrincipalUpdateRequest) (*authz.PrincipalUpdateResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	principal, err := h.principalManager.Update(req.GetId(), req.GetRoles(), attributesMap(req.GetAttributes()))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to update: %v", err.Error()))
	}

	return &authz.PrincipalUpdateResponse{
		Principal: transformer.NewPrincipal(principal).ToProto(),
	}, nil
}
