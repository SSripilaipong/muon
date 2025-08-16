package eventsource

import (
	"context"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
)

type processor struct {
	ctx      context.Context
	observer *observeSubject
	events   []CommittedEvent
}

func newProcessor(ctx context.Context, observer *observeSubject) actor.Processor[any] {
	return &processor{ctx: ctx, observer: observer}
}

func (p *processor) Process(msg any) rslt.Of[actor.Processor[any]] {
	switch msg := msg.(type) {
	case appendRequest:
		return p.processAppendRequest(msg)
	}
	return p.SameProcessor()
}

func (p *processor) SameProcessor() rslt.Of[actor.Processor[any]] {
	return rslt.Value[actor.Processor[any]](p)
}

func (p *processor) Atomic(f func(events []CommittedEvent) ([]CommittedEvent, bool)) bool {
	events, ok := f(p.events)
	if ok {
		p.events = events
	}
	return ok
}

func (p *processor) ObserverNewEvents(f func()) {
	seqBefore := p.LatestSequence()
	f()
	if len(p.events) == 0 {
		return
	}

	startIndex := seqBefore - p.events[0].Sequence() + 1
	if len(p.events) <= int(startIndex) {
		return
	}
	p.observer.Update(p.events[startIndex:])
}

func (p *processor) LatestSequence() int64 {
	if len(p.events) == 0 {
		return 0
	}
	return p.events[len(p.events)-1].Sequence()
}
