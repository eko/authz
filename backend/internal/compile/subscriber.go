package compile

import (
	"context"

	"github.com/eko/authz/backend/internal/event"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

type subscriber struct {
	logger     *slog.Logger
	compiler   Compiler
	dispatcher event.Dispatcher
}

func NewSubscriber(
	logger *slog.Logger,
	compiler Compiler,
	dispatcher event.Dispatcher,
) *subscriber {
	return &subscriber{
		logger:     logger,
		compiler:   compiler,
		dispatcher: dispatcher,
	}
}

func (s *subscriber) subscribeToPolicies(lc fx.Lifecycle) {
	policyEventChan := s.dispatcher.Subscribe(event.EventTypePolicy)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go s.handlePolicyEvents(policyEventChan)

			s.logger.Info("Compiler: subscribed to event dispatchers")

			return nil
		},
		OnStop: func(_ context.Context) error {
			close(policyEventChan)

			s.logger.Info("Compiler: subscribtion to event dispatcher stopped")

			return nil
		},
	})
}

func (s *subscriber) handlePolicyEvents(eventChan chan *event.Event) {
	for event := range eventChan {
		if err := s.compiler.CompilePolicy(event.ID); err != nil {
			s.logger.Error(
				"Compiler: unable to compile policy",
				err,
				slog.Any("policy_id", event.ID),
			)
		}
	}
}

func RunSubscriber(lc fx.Lifecycle, subscriber *subscriber) {
	subscriber.subscribeToPolicies(lc)
}
