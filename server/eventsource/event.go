package eventsource

import "github.com/SSripilaipong/muto/syntaxtree/result"

type Event interface {
	EventName() EventName
}

type EventName string

const (
	EventNameRun = "RUN"
)

type CommittedEvent struct {
	Event    Event
	sequence int64
}

func NewCommitted[E Event](event E, sequence int64) CommittedEvent {
	return CommittedEvent{
		Event:    event,
		sequence: sequence,
	}
}

func (c CommittedEvent) Sequence() int64      { return c.sequence }
func (c CommittedEvent) EventName() EventName { return c.Event.EventName() }

type RunEvent struct {
	ModuleVersion string
	Node          result.SimplifiedNode
}

func (RunEvent) EventName() EventName { return EventNameRun }

func UnsafeEventToRunEvent(e Event) RunEvent {
	return e.(RunEvent)
}
