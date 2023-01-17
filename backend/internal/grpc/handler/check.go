package handler

import (
	"context"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/pkg/authz"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Check interface {
	Check(ctx context.Context, req *authz.CheckRequest) (*authz.CheckResponse, error)
}

type check struct {
	compiledManager manager.CompiledPolicy
	logger          *slog.Logger
	dispatcher      event.Dispatcher
}

func NewCheck(
	compiledManager manager.CompiledPolicy,
	logger *slog.Logger,
	dispatcher event.Dispatcher,
) Check {
	return &check{
		compiledManager: compiledManager,
		logger:          logger,
		dispatcher:      dispatcher,
	}
}

func (h *check) Check(ctx context.Context, req *authz.CheckRequest) (*authz.CheckResponse, error) {
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
