package stats

import (
	"context"
	"regexp"
	"time"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/helper/spooler"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

type subscriber struct {
	logger            *slog.Logger
	dispatcher        event.Dispatcher
	statsManager      manager.Stats
	statsFlushDelay   time.Duration
	resourceKindRegex *regexp.Regexp
}

func NewSubscriber(
	cfg *configs.App,
	logger *slog.Logger,
	dispatcher event.Dispatcher,
	statsManager manager.Stats,
) *subscriber {
	return &subscriber{
		logger:            logger,
		dispatcher:        dispatcher,
		statsManager:      statsManager,
		statsFlushDelay:   cfg.StatsFlushDelay,
		resourceKindRegex: regexp.MustCompile(cfg.StatsResourceKindRegex),
	}
}

func (s *subscriber) subscribeToChecks(lc fx.Lifecycle) {
	checkEventChan := s.dispatcher.Subscribe(event.EventTypeCheck)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go s.handleCheckEvents(checkEventChan)

			s.logger.Info("Stats: subscribed to event dispatchers")

			return nil
		},
		OnStop: func(_ context.Context) error {
			close(checkEventChan)

			s.logger.Info("Stats: subscription to event dispatcher stopped")

			return nil
		},
	})
}

func (s *subscriber) handleCheckEvents(eventChan chan *event.Event) {
	var spooler = spooler.New(func(values []*event.Event) {
		if len(values) == 0 {
			return
		}

		var allowed, denied int64
		var timestamp int64

		for _, value := range values {
			timestamp = value.Timestamp

			checkEvent, ok := value.Data.(*event.CheckEvent)
			if !ok {
				continue
			}

			if checkEvent.IsAllowed {
				allowed++
			} else {
				denied++
			}
		}

		if err := s.statsManager.BatchAddCheck(timestamp, allowed, denied); err != nil {
			s.logger.Error("Stats: unable to add check event", err)
		}
	}, spooler.WithFlushInterval(s.statsFlushDelay))

	for eventItem := range eventChan {
		checkEvent, ok := eventItem.Data.(*event.CheckEvent)
		if !ok {
			continue
		}

		if s.resourceKindRegex.MatchString(checkEvent.ResourceKind) {
			spooler.Add(eventItem)
		}
	}
}

func RunSubscriber(lc fx.Lifecycle, subscriber *subscriber) {
	subscriber.subscribeToChecks(lc)
}
