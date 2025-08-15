package runner

import (
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"
	stResult "github.com/SSripilaipong/muto/syntaxtree/result"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	es "github.com/SSripilaipong/muon/server/eventsource"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
	"github.com/SSripilaipong/muon/server/runner/object"
)

func (s Service) Run(node stResult.SimplifiedNode) error {
	reply := make(chan error, 1)

	if err := chn.SendWithTimeout[any](s.ctrl.Ch(), runRequest{
		moduleVersion: runnerModule.VersionDefault,
		node:          node,
		reply:         reply,
	}, channelTimeout); err != nil {
		return fmt.Errorf("cannot connect to runner: %w", err)
	}
	return chn.ReceiveWithTimeout(reply, channelTimeout).Error()
}

func (p *processor) processRunRequest(msg runRequest) rslt.Of[actor.Processor[any]] {
	go func() {
		_ = chn.SendWithTimeout(msg.Reply(), p.coord.Commit([]es.Action{
			es.NewAppendAction(es.RunEvent{
				ModuleVersion: msg.ModuleVersion(),
				Node:          msg.Node(),
			}),
		}), channelTimeout)
	}()
	return rslt.Value[actor.Processor[any]](p)
}

func (p *processor) processRunEvent(event es.RunEvent, seq int64) rslt.Of[actor.Processor[any]] {
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
	return rslt.Value[actor.Processor[any]](p)
}
