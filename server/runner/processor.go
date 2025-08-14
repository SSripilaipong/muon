package runner

import (
	"context"
	"log"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	es "github.com/SSripilaipong/muon/server/eventsource"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
)

type processor struct {
	ctx              context.Context
	moduleCollection *runnerModule.Collection
	eventCtrl        *es.Controller
}

func newProcessor(ctx context.Context, moduleCollection *runnerModule.Collection, eventCtrl *es.Controller) *processor {
	return &processor{
		ctx:              ctx,
		moduleCollection: moduleCollection,
		eventCtrl:        eventCtrl,
	}
}

func (p *processor) Process(msg any) rslt.Of[actor.Processor[any]] {
	switch msg := msg.(type) {
	case runRequest:
		return p.processRunRequest(msg)
	case es.CommittedEvent:
		return p.processCommittedEvent(msg)
	default:
		log.Printf("[server.runner] unknown message type: %T", msg)
	}
	return p.SameProcessor()
}

func (p *processor) processCommittedEvent(msg es.CommittedEvent) rslt.Of[actor.Processor[any]] {
	switch msg.EventName() {
	case es.EventNameRun:
		return p.processRunEvent(es.UnsafeEventToRunEvent(msg.Event), msg.Sequence())
	default:
		log.Printf("[server.runner] unknown event name: %T", msg.EventName())
	}
	return p.SameProcessor()
}

func (p *processor) SameProcessor() rslt.Of[actor.Processor[any]] {
	return rslt.Value[actor.Processor[any]](p)
}
