package eventsource

import (
	"context"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
)

type processor struct {
	ctx         context.Context
	observer    *observeSubject
	events      []AppendedEvent
	commitUntil uint64
}

func newProcessor(ctx context.Context, observer *observeSubject) actor.Processor[any] {
	return &processor{
		ctx:      ctx,
		observer: observer,
	}
}

func (p *processor) Process(msg any) rslt.Of[actor.Processor[any]] {
	switch msg := msg.(type) {
	case appendRequest:
		return p.processAppendRequest(msg)
	case markCommitUntilRequest:
		return p.processMarkCommitUntil(msg)
	}
	return p.SameProcessor()
}

func (p *processor) SameProcessor() rslt.Of[actor.Processor[any]] {
	return rslt.Value[actor.Processor[any]](p)
}

func (p *processor) Atomic(f func(events []AppendedEvent) ([]AppendedEvent, bool)) bool {
	eventsToAppend, ok := f(p.events)
	if ok {
		p.events = append(p.events, eventsToAppend...)
	}
	return ok
}

func (p *processor) LatestSequence() uint64 {
	if len(p.events) == 0 {
		return 0
	}
	return p.events[len(p.events)-1].Sequence()
}
