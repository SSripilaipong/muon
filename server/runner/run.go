package runner

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"
	stResult "github.com/SSripilaipong/muto/syntaxtree/result"

	"github.com/SSripilaipong/muon/common/actor"
	es "github.com/SSripilaipong/muon/server/eventsource"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
	"github.com/SSripilaipong/muon/server/runner/object"
)

func (s *Service) Run(ctx context.Context, node stResult.SimplifiedNode) error {
	coord := s.coord
	if coord == nil {
		return fmt.Errorf("coordinator is not set")
	}

	err := coord.Submit(ctx, []es.Action{
		es.NewAppendAction(es.NewRunEvent(runnerModule.VersionDefault, node)),
	})
	if err != nil {
		return fmt.Errorf("cannot commit: %w", err)
	}
	return nil
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
