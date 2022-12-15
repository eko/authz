package event

type EventType string

const (
	EventTypePolicy EventType = "policy"
)

type Event struct {
	ID        string
	Timestamp int64
}
