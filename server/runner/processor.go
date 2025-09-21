package runner

import (
	"context"
	"log"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/server/coordinator"
	es "github.com/SSripilaipong/muon/server/eventsource"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
)

type processor struct {
	ctx              context.Context
	moduleCollection *runnerModule.Collection
	esStore          *es.Store
	coord            *coordinator.Controller
}

func newProcessor(ctx context.Context, moduleCollection *runnerModule.Collection, esStore *es.Store, coord *coordinator.Controller) *processor {
	return &processor{
		ctx:              ctx,
		moduleCollection: moduleCollection,
		esStore:          esStore,
		coord:            coord,
	}
}

func (p *processor) Process(msg any) rslt.Of[actor.Processor[any]] {
	switch msg := msg.(type) {
	case runRequest:
		return p.processRunRequest(msg)
	case es.AppendedEvent:
		return p.processCommittedEvent(msg)
	case localAppendRequest:
		return p.processLocalAppendRequest(msg)
	case markCommitUntilRequest:
		return p.processMarkCommitUntilRequest(msg)
	default:
		log.Printf("[server.runner] unknown message type: %T", msg)
	}
	return p.SameProcessor()
}

func (p *processor) SameProcessor() rslt.Of[actor.Processor[any]] {
	return rslt.Value[actor.Processor[any]](p)
}

func (p *processor) processCommittedEvent(msg es.AppendedEvent) rslt.Of[actor.Processor[any]] {
	switch msg.EventName() {
	case es.EventNameRun:
		return p.processRunEvent(es.UnsafeEventToRunEvent(msg.Event()), msg.Sequence())
	default:
		log.Printf("[server.runner] unknown event name: %T", msg.EventName())
	}
	return p.SameProcessor()
}

func (p *processor) processLocalAppendRequest(msg localAppendRequest) rslt.Of[actor.Processor[any]] {
	response := p.esStore.LocalAppend(p.ctx, msg.Actions())
	if err := chn.SendWithTimeout(msg.Reply(), response, channelTimeout); err != nil {
		log.Printf("[server.runner] cannot send append response: %v\n", err)
	}
	return p.SameProcessor()
}

func (p *processor) processMarkCommitUntilRequest(msg markCommitUntilRequest) rslt.Of[actor.Processor[any]] {
	err := p.esStore.MarkCommitUntil(p.ctx, msg.Sequence())
	if sendErr := chn.SendWithTimeout(msg.Reply(), err, channelTimeout); sendErr != nil {
		log.Printf("[server.runner] cannot send mark commit response: %v\n", sendErr)
	}
	return p.SameProcessor()
}
