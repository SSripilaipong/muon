package coordinator

import (
	"context"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
)

type processor struct {
	ctx   context.Context
	local Node
}

func newProcessor(ctx context.Context, local Node) *processor {
	return &processor{ctx: ctx, local: local}
}

func (p *processor) Process(msg any) rslt.Of[actor.Processor[any]] {
	switch msg := msg.(type) {
	case commitRequest:
		return p.processCommitRequest(msg)
	}
	return p.SameProcessor()
}

func (p *processor) SameProcessor() rslt.Of[actor.Processor[any]] {
	return rslt.Value[actor.Processor[any]](p)
}
