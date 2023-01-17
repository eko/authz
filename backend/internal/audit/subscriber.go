package audit

import (
	"context"
	"time"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/helper/spooler"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

type subscriber struct {
	logger          *slog.Logger
	dispatcher      event.Dispatcher
	auditManager    manager.Audit
	auditFlushDelay time.Duration
}

func NewSubscriber(
	cfg *configs.App,
	logger *slog.Logger,
	dispatcher event.Dispatcher,
	auditManager manager.Audit,
) *subscriber {
	return &subscriber{
		logger:          logger,
		dispatcher:      dispatcher,
		auditManager:    auditManager,
		auditFlushDelay: cfg.AuditFlushDelay,
	}
}

func (s *subscriber) subscribeToChecks(lc fx.Lifecycle) {
	checkEventChan := s.dispatcher.Subscribe(event.EventTypeCheck)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go s.handleCheckEvents(checkEventChan)

			s.logger.Info("Audit: subscribed to event dispatchers")

			return nil
		},
		OnStop: func(_ context.Context) error {
			close(checkEventChan)

			s.logger.Info("Audit: subscription to event dispatcher stopped")

			return nil
		},
	})
}

func (s *subscriber) handleCheckEvents(eventChan chan *event.Event) {
	var spooler = spooler.New(func(values []*event.Event) {
		if len(values) == 0 {
			return
		}

		var audits = []*model.Audit{}
		var timestamp int64

		for _, value := range values {
			timestamp = value.Timestamp

			checkEvent, ok := value.Data.(*event.CheckEvent)
			if !ok {
				continue
			}

			audit := &model.Audit{
				Date:          time.Unix(timestamp, 0),
				Principal:     checkEvent.Principal,
				ResourceKind:  checkEvent.ResourceKind,
				ResourceValue: checkEvent.ResourceValue,
				Action:        checkEvent.Action,
				IsAllowed:     checkEvent.IsAllowed,
			}

			if checkEvent.CompiledPilicy != nil {
				audit.PolicyID = checkEvent.CompiledPilicy.PolicyID
			}

			audits = append(audits, audit)
		}

		if err := s.auditManager.BatchAdd(audits); err != nil {
			s.logger.Error("Audit: unable to batch add audit events", err)
		}
	}, spooler.WithFlushInterval(s.auditFlushDelay))

	for event := range eventChan {
		spooler.Add(event)
	}
}

func RunSubscriber(lc fx.Lifecycle, subscriber *subscriber) {
	subscriber.subscribeToChecks(lc)
}
