package event

import (
	"errors"
	"sync"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/helper/time"
)

var (
	ErrEventChanCast            = errors.New("unable to cast event channel to []chan *Event")
	ErrNoSubscriberForEventType = errors.New("no subscriber for this event type")
)

type Dispatcher interface {
	Dispatch(eventType EventType, data any) error
	Subscribe(eventType EventType) chan *Event
	Unsubscribe(eventType EventType, eventChanToClose chan *Event) error
}

type dispatcher struct {
	clock         time.Clock
	subscribers   *sync.Map
	eventChanSize int
}

func NewDispatcher(
	cfg *configs.App,
	clock time.Clock,
) *dispatcher {
	return &dispatcher{
		clock:         clock,
		subscribers:   &sync.Map{},
		eventChanSize: cfg.DispatcherEventChannelSize,
	}
}

func (n *dispatcher) Dispatch(eventType EventType, data any) error {
	eventChanSlice, ok := n.subscribers.Load(eventType)
	if !ok {
		return ErrNoSubscriberForEventType
	}

	eventChans, ok := eventChanSlice.([]chan *Event)
	if !ok {
		return ErrEventChanCast
	}

	for _, eventChan := range eventChans {
		eventChan <- &Event{
			Data:      data,
			Timestamp: n.clock.Now().Unix(),
		}
	}

	return nil
}

func (n *dispatcher) Subscribe(eventType EventType) chan *Event {
	eventChan := make(chan *Event, n.eventChanSize)

	eventChanSlice, ok := n.subscribers.Load(eventType)
	if ok {
		eventChanSlice = append(eventChanSlice.([]chan *Event), eventChan)
	} else {
		eventChanSlice = []chan *Event{eventChan}
	}

	n.subscribers.Store(eventType, eventChanSlice)

	return eventChan
}

func (n *dispatcher) Unsubscribe(eventType EventType, eventChanToClose chan *Event) error {
	eventChanSlice, ok := n.subscribers.Load(eventType)
	if !ok {
		return ErrNoSubscriberForEventType
	}

	eventChans, ok := eventChanSlice.([]chan *Event)
	if !ok {
		return ErrEventChanCast
	}

	for index, eventChan := range eventChans {
		if eventChan != eventChanToClose {
			continue
		}

		eventChans = append(eventChans[:index], eventChans[index+1:]...)
		n.subscribers.LoadOrStore(eventType, eventChans)
	}

	return nil
}
