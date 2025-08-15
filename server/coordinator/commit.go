package coordinator

import (
	"context"
	"fmt"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/common/ctxs"
	es "github.com/SSripilaipong/muon/server/eventsource"
)

type commitRequest struct {
	Actions []es.Action
	Reply   chan error
}

func (c *Controller) Commit(ctx context.Context, actions []es.Action) error {
	reply := make(chan error, 1)

	var err error
	ctxs.TimeoutScope(ctx, channelTimeout, func(ctx context.Context) {
		err = chn.SendWithContext[any](ctx, c.Ch(), commitRequest{
			Actions: actions,
			Reply:   reply,
		})
	})
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
		err := p.local.Commit(p.ctx, msg.Actions)
		if err != nil {
			err = fmt.Errorf("cannot commit: %w", err)
		}
		_ = chn.SendWithTimeout(msg.Reply, err, channelTimeout)
	}()
	return p.SameProcessor()
}
