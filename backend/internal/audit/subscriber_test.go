package audit

import (
	"regexp"
	"testing"
	"time"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestNewSubscriber(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.App{
		AuditFlushDelay: 10 * time.Millisecond,
	}

	logger := slog.New(log.NewNopHandler())

	dispatcher := event.NewMockDispatcher(ctrl)

	auditManager := manager.NewMockAudit(ctrl)

	// When
	subscriberInstance := NewSubscriber(cfg, logger, dispatcher, auditManager)

	// Then
	assert := assert.New(t)

	assert.IsType(new(subscriber), subscriberInstance)

	assert.Equal(logger, subscriberInstance.logger)
	assert.Equal(dispatcher, subscriberInstance.dispatcher)
	assert.Equal(auditManager, subscriberInstance.auditManager)
	assert.Equal(cfg.AuditFlushDelay, subscriberInstance.flushDelay)
	assert.Equal(regexp.MustCompile(cfg.AuditResourceKindRegex), subscriberInstance.resourceKindRegex)
}

func TestHandleCheckEvents(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.App{
		AuditFlushDelay: 10 * time.Millisecond,
	}

	logger := slog.New(log.NewNopHandler())

	dispatcher := event.NewMockDispatcher(ctrl)

	auditManager := manager.NewMockAudit(ctrl)
	auditManager.EXPECT().BatchAdd(gomock.Len(3)).Times(1)

	subscriber := NewSubscriber(cfg, logger, dispatcher, auditManager)

	eventChan := make(chan *event.Event, 1)

	// When - Then
	go subscriber.handleCheckEvents(eventChan)

	eventChan <- &event.Event{
		Timestamp: 123456,
		Data:      &event.CheckEvent{Principal: "user1", ResourceKind: "post", ResourceValue: "1", Action: "edit", IsAllowed: true},
	}
	eventChan <- &event.Event{
		Timestamp: 123456,
		Data:      &event.CheckEvent{Principal: "user1", ResourceKind: "post", ResourceValue: "2", Action: "edit", IsAllowed: false},
	}
	eventChan <- &event.Event{
		Timestamp: 123457,
		Data:      &event.CheckEvent{Principal: "user1", ResourceKind: "post", ResourceValue: "3", Action: "delete", IsAllowed: true},
	}

	close(eventChan)

	// Wait 20ms to ensure the spool is triggered.
	<-time.After(20 * time.Millisecond)
}

func TestHandleCheckEvents_WithResourceKindRegexp(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.App{
		AuditFlushDelay:        10 * time.Millisecond,
		AuditResourceKindRegex: "^post.*",
	}

	logger := slog.New(log.NewNopHandler())

	dispatcher := event.NewMockDispatcher(ctrl)

	auditManager := manager.NewMockAudit(ctrl)
	auditManager.EXPECT().BatchAdd(gomock.Len(3)).Times(1)

	subscriber := NewSubscriber(cfg, logger, dispatcher, auditManager)

	eventChan := make(chan *event.Event, 1)

	// When - Then
	go subscriber.handleCheckEvents(eventChan)

	eventChan <- &event.Event{
		Timestamp: 123457,
		Data:      &event.CheckEvent{Principal: "user1", ResourceKind: "category", ResourceValue: "1", Action: "delete", IsAllowed: true},
	}
	eventChan <- &event.Event{
		Timestamp: 123456,
		Data:      &event.CheckEvent{Principal: "user1", ResourceKind: "post", ResourceValue: "1", Action: "edit", IsAllowed: true},
	}
	eventChan <- &event.Event{
		Timestamp: 123456,
		Data:      &event.CheckEvent{Principal: "user1", ResourceKind: "post", ResourceValue: "2", Action: "edit", IsAllowed: false},
	}
	eventChan <- &event.Event{
		Timestamp: 123457,
		Data:      &event.CheckEvent{Principal: "user1", ResourceKind: "post", ResourceValue: "3", Action: "delete", IsAllowed: true},
	}
	eventChan <- &event.Event{
		Timestamp: 123457,
		Data:      &event.CheckEvent{Principal: "user1", ResourceKind: "category", ResourceValue: "2", Action: "delete", IsAllowed: false},
	}

	close(eventChan)

	// Wait 20ms to ensure the spool is triggered.
	<-time.After(20 * time.Millisecond)
}
