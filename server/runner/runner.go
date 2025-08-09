package runner

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/chn"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
	"github.com/SSripilaipong/muon/server/runner/object"
)

type runner struct {
	ctx              context.Context
	moduleCollection *runnerModule.Collection
}

func newRunner(ctx context.Context, moduleCollection *runnerModule.Collection) runner {
	return runner{ctx: ctx, moduleCollection: moduleCollection}
}

func (r runner) Process(msg any) rslt.Of[messageProcessor] {
	switch msg := msg.(type) {
	case runMessage:
		r.ProcessRun(msg)
		return rslt.Value[messageProcessor](r)
	}
	return rslt.Value[messageProcessor](r)
}

func (r runner) ProcessRun(msg runMessage) {
	_ = chn.SendWithTimeout(msg.Reply(), func() error {
		mod, err := r.moduleCollection.Get(msg.ModuleVersion()).Return()
		if err != nil {
			return fmt.Errorf("cannot get module: %w", err)
		}

		node, ok := mod.BuildNode(msg.node.AsObject()).Return()
		if !ok {
			return fmt.Errorf("cannot build object: unknown error")
		}

		object.Spawn(r.ctx, node)
		return nil
	}(), channelTimeout)
}

type messageProcessor interface {
	Process(msg any) rslt.Of[messageProcessor]
}

func startRunner(ctx context.Context, msgBox <-chan any) <-chan struct{} {
	done := make(chan struct{})

	var p messageProcessor = newRunner(ctx, runnerModule.NewCollection())
	var err error
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgBox:
				if !ok {
					continue
				}
				p, err = p.Process(msg).Return()
				if err != nil {
					log.Printf("runner processed message and got error: %s\n", err)
					return
				}
			}
		}
	}()
	return done
}
