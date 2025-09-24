package runner

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/common/ctxs"
)

func (c *Controller) MarkCommitUntil(ctx context.Context, sequence uint64) error {
	reply := make(chan error, 1)

	err := chn.SendWithContextTimeout[any](ctx, c.Ch(), markCommitUntilRequest{
		sequence: sequence,
		reply:    reply,
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

func (p *processor) processMarkCommitUntilRequest(msg markCommitUntilRequest) rslt.Of[actor.Processor[any]] {
	err := p.esStore.MarkCommitUntil(p.ctx, msg.Sequence())
	if sendErr := chn.SendWithTimeout(msg.Reply(), err, channelTimeout); sendErr != nil {
		log.Printf("[server.runner] cannot send mark commit response: %v\n", sendErr)
	}
	return p.SameProcessor()
}
