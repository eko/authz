package handler

import (
	"context"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/pkg/authz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Check interface {
	Check(ctx context.Context, req *authz.CheckRequest) (*authz.CheckResponse, error)
}

type check struct {
	compiledManager manager.CompiledPolicy
}

func NewCheck(
	compiledManager manager.CompiledPolicy,
) Check {
	return &check{
		compiledManager: compiledManager,
	}
}

func (h *check) Check(ctx context.Context, req *authz.CheckRequest) (*authz.CheckResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var checkAnswers = make([]*authz.CheckAnswer, len(req.GetChecks()))

	for i, check := range req.GetChecks() {
		isAllowed, err := h.compiledManager.IsAllowed(check.Principal, check.ResourceKind, check.ResourceValue, check.Action)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		checkAnswers[i] = &authz.CheckAnswer{
			Principal:     check.Principal,
			ResourceKind:  check.ResourceKind,
			ResourceValue: check.ResourceValue,
			Action:        check.Action,
			IsAllowed:     isAllowed,
		}
	}

	return &authz.CheckResponse{
		Checks: checkAnswers,
	}, nil
}
