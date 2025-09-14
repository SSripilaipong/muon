package eventsource

import (
	"github.com/SSripilaipong/muto/syntaxtree/result"

	"github.com/SSripilaipong/muon/common/randutil"
)

type Event interface {
	EventName() EventName
	Hash() uint64
}

type EventName string

const (
	EventNameRun = "RUN"
)

type AppendedEvent struct {
	event    Event
	sequence uint64
}

func NewAppended[E Event](event E, sequence uint64) AppendedEvent {
	return AppendedEvent{
		event:    event,
		sequence: sequence,
	}
}

func (e AppendedEvent) Sequence() uint64     { return e.sequence }
func (e AppendedEvent) Event() Event         { return e.event }
func (e AppendedEvent) EventName() EventName { return e.Event().EventName() }
func (e AppendedEvent) Hash() uint64         { return e.Event().Hash() }

type RunEvent struct {
	ModuleVersion string
	Node          result.SimplifiedNode
	hash          uint64
}

func NewRunEvent(moduleVersion string, node result.SimplifiedNode) RunEvent {
	return RunEvent{
		ModuleVersion: moduleVersion,
		Node:          node,
		hash:          randutil.Uint64(),
	}
}

func (RunEvent) EventName() EventName { return EventNameRun }
func (e RunEvent) Hash() uint64       { return e.hash } // TODO implement hash function

func UnsafeEventToRunEvent(e Event) RunEvent {
	return e.(RunEvent)
}
