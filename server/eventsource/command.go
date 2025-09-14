package eventsource

import "github.com/SSripilaipong/go-common/optional"

type AppendAction struct {
	event            Event
	requiredSequence optional.Of[uint64]
}

func NewAppendAction(event Event, opts ...func(*AppendAction)) AppendAction {
	act := AppendAction{event: event}
	for _, opt := range opts {
		opt(&act)
	}
	return act
}

func AppendAtSequence(seq uint64) func(*AppendAction) {
	return func(a *AppendAction) {
		a.requiredSequence = optional.Value(seq)
	}
}
