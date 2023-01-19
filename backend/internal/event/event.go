package event

import "github.com/eko/authz/backend/internal/entity/model"

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

type CheckEvent struct {
	Principal      string
	ResourceKind   string
	ResourceValue  string
	Action         string
	IsAllowed      bool
	CompiledPolicy *model.CompiledPolicy
}
