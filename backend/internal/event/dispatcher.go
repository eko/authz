package event

import (
	"errors"
	"sync"

	"github.com/eko/authz/backend/internal/helper/time"
)

const (
	defaultEventChanSize uint32 = 1000
)

var (
	ErrEventChanCast            = errors.New("unable to cast event channel to []chan *Event")
	ErrNoSubscriberForEventType = errors.New("no subscriber for this event type")
)

type Dispatcher interface {
	Dispatch(eventType EventType, identifier string) error
	Subscribe(eventType EventType) chan *Event
	Unsubscribe(eventType EventType, eventChanToClose chan *Event) error
}

type dispatcher struct {
	clock         time.Clock
	subscribers   *sync.Map
	eventChanSize uint32
}

func NewDispatcher(
	clock time.Clock,
) *dispatcher {
	return &dispatcher{
		clock:         clock,
		subscribers:   &sync.Map{},
		eventChanSize: defaultEventChanSize,
	}
}

func (n *dispatcher) Dispatch(eventType EventType, identifier string) error {
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
			ID:        identifier,
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
