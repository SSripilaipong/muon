package eventsource

import (
	"fmt"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
)

type commitRequest struct {
	Actions []Action
	Reply   chan rslt.Of[CommitResult]
}

type Action any

type CommitResult struct {
	CommittedSequences []int64
}

func (c *Controller) Commit(actions []Action) rslt.Of[CommitResult] {
	reply := make(chan rslt.Of[CommitResult], 1)
	defer close(reply)

	if err := chn.SendWithTimeout[any](c.Ch(), commitRequest{
		Actions: actions,
		Reply:   reply,
	}, channelTimeout); err != nil {
		return rslt.Error[CommitResult](fmt.Errorf("cannot connect to runner: %w", err))
	}
	return rslt.Join(chn.ReceiveWithTimeout(reply, channelTimeout))
}

func (p *processor) processCommitRequest(msg commitRequest) rslt.Of[actor.Processor[any]] {
	var err error
	var committedSequences []int64

	p.Atomic(func(events []CommittedEvent) (resultEvents []CommittedEvent, ok bool) {
		events, committedSequences, err = processCommitActions(msg.Actions, p.LatestSequence(), events)
		return events, err == nil
	})

	_ = chn.SendWithTimeout(msg.Reply, func() rslt.Of[CommitResult] {
		if err != nil {
			return rslt.Error[CommitResult](err)
		}
		return rslt.Value(CommitResult{CommittedSequences: committedSequences})
	}(), channelTimeout)
	fmt.Println("current events:", p.events)
	return p.SameProcessor()
}

func processCommitActions(actions []Action, seq int64, events []CommittedEvent) ([]CommittedEvent, []int64, error) {
	var cs []int64
	for _, action := range actions {
		switch action := action.(type) {
		case AppendAction:
			rs := action.requiredSequence
			if rs.IsNotEmpty() && rs.Value() != seq {
				return nil, nil, fmt.Errorf("sequence requirement violation")
			}
			seq++
			events = append(events, NewCommitted(action.event, seq))
			cs = append(cs, seq)
		default:
			return nil, nil, fmt.Errorf("unknown action %T", action)
		}
	}
	return events, cs, nil
}
