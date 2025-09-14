package eventsource

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/common/ctxs"
	"github.com/SSripilaipong/muon/common/msgutil"
	"github.com/SSripilaipong/muon/common/slc"
)

type appendRequest struct {
	msgutil.ReplyMixin[rslt.Of[AppendResponse]]
	Actions []Action
}

type Action any

type AppendResponse struct {
	LatestCommittedSequence uint64
	ChainedEvents           []ChainedEvent
}

type ChainedEvent struct {
	Event            AppendedEvent
	PreviousSequence uint64
	PreviousHash     uint64
}

func (c *Controller) LocalAppend(ctx context.Context, actions []Action) rslt.Of[AppendResponse] {
	reply := make(chan rslt.Of[AppendResponse], 1)

	err := chn.SendWithContextTimeout[any](ctx, c.Ch(), appendRequest{
		Actions:    actions,
		ReplyMixin: msgutil.NewReplyMixin(reply, channelTimeout),
	}, channelTimeout)
	if err != nil {
		return rslt.Error[AppendResponse](fmt.Errorf("cannot connect to event source: %w", err))
	}

	var response rslt.Of[AppendResponse]
	ctxs.TimeoutScope(ctx, channelTimeout, func(ctx context.Context) {
		response = rslt.Join(chn.ReceiveWithContext(ctx, reply))
	})
	return response
}

func (p *processor) processAppendRequest(msg appendRequest) rslt.Of[actor.Processor[any]] {
	var eventsToAppend []AppendedEvent
	var previousHash uint64
	var appendErr error

	p.Atomic(func(events []AppendedEvent) (resultEvents []AppendedEvent, ok bool) { // actually the events param should be the reduced current state, but it can't be implemented now
		previousHash = slc.LastDefaultZero(events).Hash()

		eventsToAppend, appendErr = processAppendActions(msg.Actions, p.LatestSequence())
		return eventsToAppend, appendErr == nil
	})
	appendedEvents := eventsToAppend

	go respondAppend(msg.Reply, appendErr, appendedEvents, previousHash, p.LatestSequence())

	log.Println("DEBUG current events:", p.events)
	return p.SameProcessor()
}

func respondAppend(respond func(x rslt.Of[AppendResponse]) error, appendErr error, appendedEvents []AppendedEvent, latestHash uint64, latestSequence uint64) { // TODO implement
	_ = respond(func() rslt.Of[AppendResponse] {
		if appendErr != nil {
			return rslt.Error[AppendResponse](appendErr)
		}

		var chainedEvents []ChainedEvent
		previousHash, previousSequence := latestHash, latestSequence
		for _, event := range appendedEvents {
			chainedEvents = append(chainedEvents, ChainedEvent{
				Event:            event,
				PreviousSequence: previousSequence,
				PreviousHash:     previousHash,
			})
			previousHash, previousSequence = event.Hash(), event.Sequence()
		}
		return rslt.Value(AppendResponse{
			LatestCommittedSequence: latestSequence,
			ChainedEvents:           chainedEvents,
		})
	}())
}

func processAppendActions(actions []Action, previousSeq uint64) ([]AppendedEvent, error) {
	var cs []uint64
	var eventsToAppend []AppendedEvent
	seq := previousSeq
	for _, action := range actions {
		switch action := action.(type) {
		case AppendAction:
			rs := action.requiredSequence
			if rs.IsNotEmpty() && rs.Value() != seq {
				return nil, fmt.Errorf("sequence requirement violation")
			}
			seq++
			eventsToAppend = append(eventsToAppend, NewAppended(action.event, seq))
			cs = append(cs, seq)
		default:
			return nil, fmt.Errorf("unknown action %T", action)
		}
	}
	return eventsToAppend, nil
}
