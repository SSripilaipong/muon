package coordinator

import (
	"context"
	"fmt"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/common/ctxs"
	"github.com/SSripilaipong/muon/common/prl"
	es "github.com/SSripilaipong/muon/server/eventsource"
)

type commitRequest struct {
	Actions []es.Action
	Reply   chan error
}

func (c *Controller) Commit(ctx context.Context, actions []es.Action) error {
	reply := make(chan error, 1)

	err := chn.SendWithContextTimeout[any](ctx, c.Ch(), commitRequest{
		Actions: actions,
		Reply:   reply,
	}, channelTimeout)
	if err != nil {
		return fmt.Errorf("cannot connect to coordinator: %w", err)
	}

	var response error
	ctxs.TimeoutScope(ctx, channelTimeout, func(ctx context.Context) {
		response = rslt.JoinError(chn.ReceiveWithContext(ctx, reply))
	})
	return response
}

func (p *processor) processCommitRequest(msg commitRequest) rslt.Of[actor.Processor[any]] {
	go func() {
		_ = chn.SendWithTimeout(msg.Reply, func() error {
			if err := p.localNode.LocalCommit(p.ctx, msg.Actions); err != nil {
				return fmt.Errorf("cannot commit to local node: %w", err)
			}

			responses := prl.Collect(p.ctx, commitsFromNodes(p.ctx, p.nodes, msg.Actions)...)
			ok, _ := guaranteeQuorumSuccess(p.ctx, len(p.nodes), responses)
			if !ok {
				return fmt.Errorf("cannot guarantee quorum commit")
			}

			return nil
		}(), channelTimeout)
	}()
	return p.SameProcessor()
}

func commitsFromNodes(ctx context.Context, nodes []Node, actions []es.Action) (result []func() error) {
	for _, node := range nodes {
		result = append(result, func() error {
			return node.Commit(ctx, actions)
		})
	}
	return result
}

func guaranteeQuorumSuccess(ctx context.Context, n int, sources <-chan error) (bool, []error) {
	defer func() { go func() { chn.Drain(sources) }() }()

	var errs []error
loop:
	for {
		select {
		case err, isOpen := <-sources:
			if !isOpen {
				break loop
			}

			if err != nil {
				errs = append(errs, err)
			}
			if len(errs) > n/2 {
				return false, errs
			}
		case <-ctx.Done():
			return false, errs
		}
	}
	return true, errs
}
