package eventsource

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/common/ctxs"
	"github.com/SSripilaipong/muon/common/fn"
	"github.com/SSripilaipong/muon/common/msgutil"
)

type markCommitUntilRequest struct {
	Sequence uint64
	msgutil.ReplyMixin[error]
}

func (c *Controller) MarkCommitUntil(ctx context.Context, sequence uint64) error {
	reply := make(chan error, 1)

	err := chn.SendWithContextTimeout[any](ctx, c.Ch(), markCommitUntilRequest{
		Sequence:   sequence,
		ReplyMixin: msgutil.NewReplyMixin(reply, channelTimeout),
	}, channelTimeout)
	if err != nil {
		return fmt.Errorf("cannot connect to event source: %w", err)
	}

	var response error
	ctxs.TimeoutScope(ctx, channelTimeout, func(ctx context.Context) {
		response = rslt.Transform(fn.Id[error], fn.Id)(chn.ReceiveWithContext(ctx, reply))
	})
	return response
}

func (p *processor) processMarkCommitUntil(msg markCommitUntilRequest) rslt.Of[actor.Processor[any]] {
	previousCommitUntil := p.commitUntil

	_ = msg.Reply(func() error {
		if msg.Sequence < p.commitUntil {
			return fmt.Errorf("decreasing commit is not allowed")
		}
		p.commitUntil = msg.Sequence
		return nil
	}())

	commitStartIndex, err := seekToSequence(p.events, previousCommitUntil+1).Return()
	if err != nil {
		log.Println("[server.eventsource] error while seeking sequence:", err)
		return p.SameProcessor()
	}
	p.observer.Update(p.events[commitStartIndex:])

	log.Println("DEBUG commit until:", p.commitUntil)
	return p.SameProcessor()
}

func seekToSequence(events []AppendedEvent, seq uint64) rslt.Of[uint64] {
	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		if event.Sequence() == seq {
			return rslt.Value(uint64(i))
		} else if event.Sequence() < seq {
			return rslt.Error[uint64](fmt.Errorf("out of sequence"))
		}
	}
	return rslt.Error[uint64](fmt.Errorf("sequence not found"))
}
