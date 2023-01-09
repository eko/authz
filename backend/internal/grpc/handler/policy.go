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

type Policy interface {
	PolicyCreate(ctx context.Context, req *authz.PolicyCreateRequest) (*authz.PolicyCreateResponse, error)
	PolicyDelete(ctx context.Context, req *authz.PolicyDeleteRequest) (*authz.PolicyDeleteResponse, error)
	PolicyGet(ctx context.Context, req *authz.PolicyGetRequest) (*authz.PolicyGetResponse, error)
	PolicyUpdate(ctx context.Context, req *authz.PolicyUpdateRequest) (*authz.PolicyUpdateResponse, error)
}

type policy struct {
	policyManager manager.Policy
}

func NewPolicy(
	policyManager manager.Policy,
) Policy {
	return &policy{
		policyManager: policyManager,
	}
}

func (h *policy) PolicyCreate(ctx context.Context, req *authz.PolicyCreateRequest) (*authz.PolicyCreateResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	policy, err := h.policyManager.Create(req.GetId(), req.GetResources(), req.GetActions(), req.GetAttributeRules())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to create: %v", err.Error()))
	}

	return &authz.PolicyCreateResponse{
		Policy: transformer.NewPolicy(policy).ToProto(),
	}, nil
}

func (h *policy) PolicyDelete(ctx context.Context, req *authz.PolicyDeleteRequest) (*authz.PolicyDeleteResponse, error) {
	err := h.policyManager.Delete(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to delete: %v", err.Error()))
	}

	return &authz.PolicyDeleteResponse{
		Success: true,
	}, nil
}

func (h *policy) PolicyGet(ctx context.Context, req *authz.PolicyGetRequest) (*authz.PolicyGetResponse, error) {
	policy, err := h.policyManager.GetRepository().Get(req.GetId(), repository.WithPreloads("Attributes", "Roles"))
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to retrieve: %v", err.Error()))
	}

	return &authz.PolicyGetResponse{
		Policy: transformer.NewPolicy(policy).ToProto(),
	}, nil
}

func (h *policy) PolicyUpdate(ctx context.Context, req *authz.PolicyUpdateRequest) (*authz.PolicyUpdateResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	policy, err := h.policyManager.Update(req.GetId(), req.GetResources(), req.GetActions(), req.GetAttributeRules())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to update: %v", err.Error()))
	}

	return &authz.PolicyUpdateResponse{
		Policy: transformer.NewPolicy(policy).ToProto(),
	}, nil
}
