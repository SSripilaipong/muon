package runner

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/common/ctxs"
	es "github.com/SSripilaipong/muon/server/eventsource"
)

func (c *Controller) LocalAppend(ctx context.Context, actions []es.Action) rslt.Of[es.AppendResponse] {
	reply := make(chan rslt.Of[es.AppendResponse], 1)

	err := chn.SendWithContextTimeout[any](ctx, c.Ch(), localAppendRequest{
		actions: actions,
		reply:   reply,
	}, channelTimeout)
	if err != nil {
		return rslt.Error[es.AppendResponse](fmt.Errorf("cannot connect to runner: %w", err))
	}

	var response rslt.Of[es.AppendResponse]
	ctxs.TimeoutScope(ctx, channelTimeout, func(ctx context.Context) {
		response = rslt.Join(chn.ReceiveWithContext(ctx, reply))
	})
	return response
}

func (p *processor) processLocalAppendRequest(msg localAppendRequest) rslt.Of[actor.Processor[any]] {
	response := p.esStore.LocalAppend(p.ctx, msg.Actions())
	if err := chn.SendWithTimeout(msg.Reply(), response, channelTimeout); err != nil {
		log.Printf("[server.runner] cannot send append response: %v\n", err)
	}
	return p.SameProcessor()
}
