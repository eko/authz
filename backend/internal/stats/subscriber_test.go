package stats

import (
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
		StatsFlushDelay: 10 * time.Millisecond,
	}

	logger := slog.New(log.NewNopHandler())

	dispatcher := event.NewMockDispatcher(ctrl)

	statsManager := manager.NewMockStats(ctrl)

	// When
	subscriberInstance := NewSubscriber(cfg, logger, dispatcher, statsManager)

	// Then
	assert := assert.New(t)

	assert.IsType(new(subscriber), subscriberInstance)

	assert.Equal(logger, subscriberInstance.logger)
	assert.Equal(dispatcher, subscriberInstance.dispatcher)
	assert.Equal(statsManager, subscriberInstance.statsManager)
	assert.Equal(cfg.StatsFlushDelay, subscriberInstance.statsFlushDelay)
}

func TestHandleCheckEvents(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.App{
		StatsFlushDelay: 10 * time.Millisecond,
	}

	logger := slog.New(log.NewNopHandler())

	dispatcher := event.NewMockDispatcher(ctrl)

	statsManager := manager.NewMockStats(ctrl)
	statsManager.EXPECT().BatchAddCheck(int64(123457), int64(2), int64(1)).Times(1)

	subscriber := NewSubscriber(cfg, logger, dispatcher, statsManager)

	eventChan := make(chan *event.Event, 1)

	// When - Then
	go subscriber.handleCheckEvents(eventChan)

	eventChan <- &event.Event{Timestamp: 123456, Data: true}
	eventChan <- &event.Event{Timestamp: 123456, Data: false}
	eventChan <- &event.Event{Timestamp: 123457, Data: true}

	close(eventChan)

	// Wait 20ms to ensure the spool is triggered.
	<-time.After(20 * time.Millisecond)
}
