package metric

import (
	"context"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/event"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

type subscriber struct {
	enabled    bool
	logger     *slog.Logger
	dispatcher event.Dispatcher
	observer   Observer
}

func NewSubscriber(
	cfg *configs.App,
	logger *slog.Logger,
	dispatcher event.Dispatcher,
	observer Observer,
) *subscriber {
	return &subscriber{
		enabled:    cfg.MetricsEnabled,
		logger:     logger,
		dispatcher: dispatcher,
		observer:   observer,
	}
}

func (s *subscriber) subscribe(lc fx.Lifecycle) {
	if !s.enabled {
		return
	}

	checkEventChan := s.dispatcher.Subscribe(event.EventTypeCheck)
	policyEventChan := s.dispatcher.Subscribe(event.EventTypePolicy)
	principalEventChan := s.dispatcher.Subscribe(event.EventTypePrincipal)
	resourceEventChan := s.dispatcher.Subscribe(event.EventTypeResource)
	roleEventChan := s.dispatcher.Subscribe(event.EventTypeRole)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go s.handleCheckEvents(checkEventChan)

			go s.handleItemEvents(policyEventChan, "policy")
			go s.handleItemEvents(principalEventChan, "principal")
			go s.handleItemEvents(resourceEventChan, "resource")
			go s.handleItemEvents(roleEventChan, "role")

			s.logger.Info("Metric: subscribed to event dispatchers")

			return nil
		},
		OnStop: func(_ context.Context) error {
			close(checkEventChan)
			close(policyEventChan)

			s.logger.Info("Metric: subscription to event dispatcher stopped")

			return nil
		},
	})
}

func (s *subscriber) handleCheckEvents(eventChan chan *event.Event) {
	for eventItem := range eventChan {
		if !s.enabled {
			continue
		}

		checkEvent, ok := eventItem.Data.(*event.CheckEvent)
		if !ok {
			continue
		}

		s.observer.ObserveCheckCounter(checkEvent.ResourceKind, checkEvent.IsAllowed)
	}
}

func (s *subscriber) handleItemEvents(eventChan chan *event.Event, itemType string) {
	for eventItem := range eventChan {
		if !s.enabled {
			continue
		}

		itemEvent, ok := eventItem.Data.(*event.ItemEvent)
		if !ok {
			continue
		}

		s.observer.ObserveItemCreatedCounter(itemType, string(itemEvent.Action))
	}
}

func RunSubscriber(lc fx.Lifecycle, subscriber *subscriber) {
	subscriber.subscribe(lc)
}
