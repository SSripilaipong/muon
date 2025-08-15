package eventsource

import (
	"fmt"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
)

type commitRequest struct {
	Actions []Action
	Reply   chan error
}

type Action any

func (c *Controller) Commit(actions []Action) error {
	reply := make(chan error, 1)

	if err := chn.SendWithTimeout[any](c.Ch(), commitRequest{
		Actions: actions,
		Reply:   reply,
	}, channelTimeout); err != nil {
		return fmt.Errorf("cannot connect to runner: %w", err)
	}
	return chn.ReceiveWithTimeout(reply, channelTimeout).Error()
}

func (p *processor) processCommitRequest(msg commitRequest) rslt.Of[actor.Processor[any]] {
	var err error

	p.ObserverNewEvents(func() {
		p.Atomic(func(events []CommittedEvent) (resultEvents []CommittedEvent, ok bool) {
			events, err = processCommitActions(msg.Actions, p.LatestSequence(), events)
			return events, err == nil
		})
	})

	go func() {
		_ = chn.SendWithTimeout(msg.Reply, err, channelTimeout)
	}()
	fmt.Println("current events:", p.events)
	return p.SameProcessor()
}

func processCommitActions(actions []Action, seq int64, events []CommittedEvent) ([]CommittedEvent, error) {
	var cs []int64
	for _, action := range actions {
		switch action := action.(type) {
		case AppendAction:
			rs := action.requiredSequence
			if rs.IsNotEmpty() && rs.Value() != seq {
				return nil, fmt.Errorf("sequence requirement violation")
			}
			seq++
			events = append(events, NewCommitted(action.event, seq))
			cs = append(cs, seq)
		default:
			return nil, fmt.Errorf("unknown action %T", action)
		}
	}
	return events, nil
}
