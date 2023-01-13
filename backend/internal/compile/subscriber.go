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
	principalEventChan := s.dispatcher.Subscribe(event.EventTypePrincipal)
	resourceEventChan := s.dispatcher.Subscribe(event.EventTypeResource)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go s.handlePolicyEvents(policyEventChan)
			go s.handlePrincipalEvents(principalEventChan)
			go s.handleResourceEvents(resourceEventChan)

			s.logger.Info("Compiler: subscribed to event dispatchers")

			return nil
		},
		OnStop: func(_ context.Context) error {
			close(policyEventChan)
			close(principalEventChan)
			close(resourceEventChan)

			s.logger.Info("Compiler: subscription to event dispatcher stopped")

			return nil
		},
	})
}

func (s *subscriber) handlePolicyEvents(eventChan chan *event.Event) {
	for event := range eventChan {
		if err := s.compiler.CompilePolicy(event.Data.(string)); err != nil {
			s.logger.Warn(
				"Compiler: unable to compile policy",
				err,
				slog.Any("policy_id", event.Data.(string)),
			)
		}
	}
}

func (s *subscriber) handleResourceEvents(eventChan chan *event.Event) {
	for event := range eventChan {
		if err := s.compiler.CompileResource(event.Data.(string)); err != nil {
			s.logger.Warn(
				"Compiler: unable to compile resource",
				err,
				slog.Any("policy_id", event.Data.(string)),
			)
		}
	}
}

func (s *subscriber) handlePrincipalEvents(eventChan chan *event.Event) {
	for event := range eventChan {
		if err := s.compiler.CompilePrincipal(event.Data.(string)); err != nil {
			s.logger.Warn(
				"Compiler: unable to compile principal",
				err,
				slog.Any("policy_id", event.Data.(string)),
			)
		}
	}
}

func RunSubscriber(lc fx.Lifecycle, subscriber *subscriber) {
	subscriber.subscribeToPolicies(lc)
}
