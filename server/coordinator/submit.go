package coordinator

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/common/ctxs"
	"github.com/SSripilaipong/muon/common/errorutil"
	"github.com/SSripilaipong/muon/common/msgutil"
	"github.com/SSripilaipong/muon/common/prl"
	"github.com/SSripilaipong/muon/common/rsltutil"
	es "github.com/SSripilaipong/muon/server/eventsource"
)

type submitRequest struct {
	msgutil.ReplyMixin[error]
	Actions []es.Action
}

func (c *Controller) Submit(ctx context.Context, actions []es.Action) error {
	reply := make(chan error, 1)

	err := chn.SendWithContextTimeout[any](ctx, c.Ch(), submitRequest{
		Actions:    actions,
		ReplyMixin: msgutil.NewReplyMixin(reply, channelTimeout),
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

func (p *processor) processSubmitRequest(msg submitRequest) rslt.Of[actor.Processor[any]] {
	go submit(p.ctx, p.localNode, p.nodes, p.NumberOfOnlineNodes(), msg)
	return p.SameProcessor()
}

func submit(ctx context.Context, localNode LocalNode, otherNodes []Node, numberOfOnlineNodes int, msg submitRequest) {
	wrapAppendError := rsltutil.WrapError[es.AppendResponse](errorutil.Wrapf("cannot append to local node: %w"))
	appendResponse, appendError := wrapAppendError(localNode.LocalAppend(ctx, msg.Actions)).Return()

	_ = msg.Reply(appendError)
	if appendError != nil {
		return
	}

	// TODO make all code below as an effect from calling LocalAppend() to localNode <- after that, what is the purpose of coordinator then?
	responses := prl.Collect(ctx, appendingResultFromNodes(ctx, otherNodes, msg.Actions)...)
	ok, _ := guaranteeAtLeastHalfSuccess(ctx, numberOfOnlineNodes, responses)
	if !ok {
		log.Println("[server.coordinator] appending to other nodes failed")
		return
	}

	if err := localNode.MarkCommitUntil(ctx, appendResponse.LatestCommittedSequence); err != nil {
		log.Println("[server.coordinator] fail to mark commit until to local node")
		return
	}
}

func appendingResultFromNodes(ctx context.Context, nodes []Node, actions []es.Action) (result []func() error) {
	for _, node := range nodes {
		result = append(result, func() error {
			return node.ForceAppend(ctx, actions).Error()
		})
	}
	return result
}

func guaranteeAtLeastHalfSuccess(ctx context.Context, n int, sources <-chan error) (bool, []error) {
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
			if float64(len(errs)) >= float64(n)/2 {
				return false, errs
			}
		case <-ctx.Done():
			return false, errs
		}
	}
	return true, errs
}
