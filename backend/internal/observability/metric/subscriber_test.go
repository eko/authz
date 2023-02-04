package metric

import (
	"testing"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/model"
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
		MetricsEnabled: true,
	}

	logger := slog.New(log.NewNopHandler())

	dispatcher := event.NewMockDispatcher(ctrl)

	observer := NewMockObserver(ctrl)

	// When
	subscriberInstance := NewSubscriber(cfg, logger, dispatcher, observer)

	// Then
	assert := assert.New(t)

	assert.IsType(new(subscriber), subscriberInstance)

	assert.Equal(cfg.MetricsEnabled, subscriberInstance.enabled)
	assert.Equal(logger, subscriberInstance.logger)
	assert.Equal(dispatcher, subscriberInstance.dispatcher)
	assert.Equal(observer, subscriberInstance.observer)
}

func TestHandleCheckEvents_WhenEnabled(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.App{
		MetricsEnabled: true,
	}

	logger := slog.New(log.NewNopHandler())

	dispatcher := event.NewMockDispatcher(ctrl)

	observer := NewMockObserver(ctrl)
	observer.EXPECT().ObserveCheckCounter("post", true).Times(2)
	observer.EXPECT().ObserveCheckCounter("post", false).Times(1)

	subscriber := NewSubscriber(cfg, logger, dispatcher, observer)

	eventChan := make(chan *event.Event)

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
}

func TestHandleCheckEvents_WhenNotEnabled(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.App{
		MetricsEnabled: false,
	}

	logger := slog.New(log.NewNopHandler())

	dispatcher := event.NewMockDispatcher(ctrl)

	observer := NewMockObserver(ctrl)

	subscriber := NewSubscriber(cfg, logger, dispatcher, observer)

	eventChan := make(chan *event.Event)

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
}

func TestHandleItemEvents_WhenEnabled(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.App{
		MetricsEnabled: true,
	}

	logger := slog.New(log.NewNopHandler())

	dispatcher := event.NewMockDispatcher(ctrl)

	observer := NewMockObserver(ctrl)
	observer.EXPECT().ObserveItemCreatedCounter("resource", "create").Times(2)
	observer.EXPECT().ObserveItemCreatedCounter("resource", "update").Times(1)

	subscriber := NewSubscriber(cfg, logger, dispatcher, observer)

	eventChan := make(chan *event.Event)

	// When - Then
	go subscriber.handleItemEvents(eventChan, "resource")

	eventChan <- &event.Event{
		Timestamp: 123456,
		Data:      &event.ItemEvent{Action: event.ItemActionCreate, Data: &model.Resource{ID: "4"}},
	}
	eventChan <- &event.Event{
		Timestamp: 123456,
		Data:      &event.ItemEvent{Action: event.ItemActionUpdate, Data: &model.Resource{ID: "4"}},
	}
	eventChan <- &event.Event{
		Timestamp: 123456,
		Data:      &event.ItemEvent{Action: event.ItemActionCreate, Data: &model.Resource{ID: "5"}},
	}

	close(eventChan)
}

func TestHandleItemEvents_WhenNotEnabled(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	cfg := &configs.App{
		MetricsEnabled: false,
	}

	logger := slog.New(log.NewNopHandler())

	dispatcher := event.NewMockDispatcher(ctrl)

	observer := NewMockObserver(ctrl)

	subscriber := NewSubscriber(cfg, logger, dispatcher, observer)

	eventChan := make(chan *event.Event)

	// When - Then
	go subscriber.handleItemEvents(eventChan, "resource")

	eventChan <- &event.Event{
		Timestamp: 123456,
		Data:      &event.ItemEvent{Action: event.ItemActionCreate, Data: &model.Resource{ID: "4"}},
	}
	eventChan <- &event.Event{
		Timestamp: 123456,
		Data:      &event.ItemEvent{Action: event.ItemActionUpdate, Data: &model.Resource{ID: "4"}},
	}
	eventChan <- &event.Event{
		Timestamp: 123456,
		Data:      &event.ItemEvent{Action: event.ItemActionCreate, Data: &model.Resource{ID: "5"}},
	}

	close(eventChan)
}
