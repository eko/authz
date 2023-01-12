package event

type EventType string

const (
	EventTypeCheck     EventType = "check"
	EventTypePolicy    EventType = "policy"
	EventTypePrincipal EventType = "principal"
	EventTypeResource  EventType = "resource"
)

type Event struct {
	Data      any
	Timestamp int64
}
