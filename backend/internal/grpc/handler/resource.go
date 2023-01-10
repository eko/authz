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

type Resource interface {
	ResourceCreate(ctx context.Context, req *authz.ResourceCreateRequest) (*authz.ResourceCreateResponse, error)
	ResourceDelete(ctx context.Context, req *authz.ResourceDeleteRequest) (*authz.ResourceDeleteResponse, error)
	ResourceGet(ctx context.Context, req *authz.ResourceGetRequest) (*authz.ResourceGetResponse, error)
	ResourceUpdate(ctx context.Context, req *authz.ResourceUpdateRequest) (*authz.ResourceUpdateResponse, error)
}

type resource struct {
	resourceManager manager.Resource
}

func NewResource(
	resourceManager manager.Resource,
) Resource {
	return &resource{
		resourceManager: resourceManager,
	}
}

func (h *resource) ResourceCreate(ctx context.Context, req *authz.ResourceCreateRequest) (*authz.ResourceCreateResponse, error) {
	if !validator.ValidateSlugFromString(req.GetId()) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("identifier must be a slug, found: %s", req.GetId()))
	}

	resource, err := h.resourceManager.Create(req.GetId(), req.GetKind(), req.GetValue(), attributesMap(req.GetAttributes()))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to create: %v", err.Error()))
	}

	return &authz.ResourceCreateResponse{
		Resource: transformer.NewResource(resource).ToProto(),
	}, nil
}

func (h *resource) ResourceDelete(ctx context.Context, req *authz.ResourceDeleteRequest) (*authz.ResourceDeleteResponse, error) {
	err := h.resourceManager.Delete(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to delete: %v", err.Error()))
	}

	return &authz.ResourceDeleteResponse{
		Success: true,
	}, nil
}

func (h *resource) ResourceGet(ctx context.Context, req *authz.ResourceGetRequest) (*authz.ResourceGetResponse, error) {
	resource, err := h.resourceManager.GetRepository().Get(req.GetId(), repository.WithPreloads("Attributes"))
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to retrieve: %v", err.Error()))
	}

	return &authz.ResourceGetResponse{
		Resource: transformer.NewResource(resource).ToProto(),
	}, nil
}

func (h *resource) ResourceUpdate(ctx context.Context, req *authz.ResourceUpdateRequest) (*authz.ResourceUpdateResponse, error) {
	resource, err := h.resourceManager.Update(req.GetId(), req.GetKind(), req.GetValue(), attributesMap(req.GetAttributes()))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to update: %v", err.Error()))
	}

	return &authz.ResourceUpdateResponse{
		Resource: transformer.NewResource(resource).ToProto(),
	}, nil
}
