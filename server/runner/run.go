package runner

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"
	stResult "github.com/SSripilaipong/muto/syntaxtree/result"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/common/ctxs"
	es "github.com/SSripilaipong/muon/server/eventsource"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
	"github.com/SSripilaipong/muon/server/runner/object"
)

func (s Service) Run(ctx context.Context, node stResult.SimplifiedNode) error {
	reply := make(chan error, 1)

	err := chn.SendWithContextTimeout[any](ctx, s.ctrl.Ch(), runRequest{
		moduleVersion: runnerModule.VersionDefault,
		node:          node,
		reply:         reply,
	}, channelTimeout)
	if err != nil {
		return fmt.Errorf("cannot connect to runner: %w", err)
	}

	var response error
	ctxs.TimeoutScope(ctx, channelTimeout, func(ctx context.Context) {
		response = rslt.JoinError(chn.ReceiveWithContext(ctx, reply))
	})
	return response
}

func (p *processor) processRunRequest(msg runRequest) rslt.Of[actor.Processor[any]] {
	go func() {
		coord := p.coord
		if coord == nil {
			if err := chn.SendWithTimeout(msg.Reply(), fmt.Errorf("coordinator is not set"), channelTimeout); err != nil {
				log.Printf("[server.runner] cannot send run error: %v\n", err)
			}
			return
		}

		err := coord.Submit(p.ctx, []es.Action{
			es.NewAppendAction(es.NewRunEvent(msg.ModuleVersion(), msg.Node())),
		})
		if err != nil {
			err = fmt.Errorf("cannot commit: %w", err)
		}
		if sendErr := chn.SendWithTimeout(msg.Reply(), err, channelTimeout); sendErr != nil {
			log.Printf("[server.runner] cannot send run response: %v\n", sendErr)
		}
	}()
	return p.SameProcessor()
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

func (p *processor) processRunEvent(event es.RunEvent, _ uint64) rslt.Of[actor.Processor[any]] {
	if err := func() error {
		mod, err := p.moduleCollection.Get(event.ModuleVersion).Return()
		if err != nil {
			return fmt.Errorf("cannot get module: %w", err)
		}

		node, ok := mod.BuildNode(event.Node.AsObject()).Return()
		if !ok {
			return fmt.Errorf("cannot build object: unknown error")
		}

		object.Spawn(p.ctx, node)
		return nil
	}(); err != nil {
		log.Printf("[server.runner] fail to process run event: %v\n", err)
	}
	return p.SameProcessor()
}
