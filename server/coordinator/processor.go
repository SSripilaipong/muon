package coordinator

import (
	"context"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
)

type processor struct {
	ctx       context.Context
	localNode LocalNode
	nodes     []Node
}

func newProcessor(ctx context.Context, local LocalNode) *processor {
	return &processor{ctx: ctx, localNode: local, nodes: []Node{}}
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
