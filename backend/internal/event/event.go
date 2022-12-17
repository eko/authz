package event

type EventType string

const (
	EventTypePolicy    EventType = "policy"
	EventTypePrincipal EventType = "principal"
	EventTypeResource  EventType = "resource"
)

type Event struct {
	ID        string
	Timestamp int64
}
