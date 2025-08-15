package coordinator

import (
	"fmt"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	es "github.com/SSripilaipong/muon/server/eventsource"
)

type commitRequest struct {
	Actions []es.Action
	Reply   chan error
}

func (c *Controller) Commit(actions []es.Action) error {
	reply := make(chan error, 1)

	if err := chn.SendWithTimeout[any](c.Ch(), commitRequest{
		Actions: actions,
		Reply:   reply,
	}, channelTimeout); err != nil {
		return fmt.Errorf("cannot connect to coordinator: %w", err)
	}
	return chn.ReceiveWithTimeout(reply, channelTimeout).Error()
}

func (p *processor) processCommitRequest(msg commitRequest) rslt.Of[actor.Processor[any]] {
	go func() {
		_ = chn.SendWithTimeout(msg.Reply, p.local.Commit(msg.Actions), channelTimeout)
	}()
	return p.SameProcessor()
}
