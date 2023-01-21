package compile

import (
	"context"

	"github.com/eko/authz/backend/internal/entity/model"
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
	for eventItem := range eventChan {
		itemEvent, ok := eventItem.Data.(*event.ItemEvent)
		if !ok {
			continue
		}

		policy := itemEvent.Data.(*model.Policy)

		if err := s.compiler.CompilePolicy(policy); err != nil {
			s.logger.Warn(
				"Compiler: unable to compile policy",
				err,
				slog.Any("policy_id", policy.ID),
			)
		}
	}
}

func (s *subscriber) handleResourceEvents(eventChan chan *event.Event) {
	for eventItem := range eventChan {
		itemEvent, ok := eventItem.Data.(*event.ItemEvent)
		if !ok {
			continue
		}

		resource := itemEvent.Data.(*model.Resource)

		if err := s.compiler.CompileResource(resource); err != nil {
			s.logger.Warn(
				"Compiler: unable to compile resource",
				err,
				slog.Any("policy_id", resource.ID),
			)
		}
	}
}

func (s *subscriber) handlePrincipalEvents(eventChan chan *event.Event) {
	for eventItem := range eventChan {
		itemEvent, ok := eventItem.Data.(*event.ItemEvent)
		if !ok {
			continue
		}

		principal := itemEvent.Data.(*model.Principal)

		if err := s.compiler.CompilePrincipal(principal); err != nil {
			s.logger.Warn(
				"Compiler: unable to compile principal",
				err,
				slog.Any("policy_id", principal.ID),
			)
		}
	}
}

func RunSubscriber(lc fx.Lifecycle, subscriber *subscriber) {
	subscriber.subscribeToPolicies(lc)
}
