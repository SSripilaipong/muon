package runner

import (
	"context"
	"fmt"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/common/ctxs"
	"github.com/SSripilaipong/muon/server/coordinator"
	es "github.com/SSripilaipong/muon/server/eventsource"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
)

type Controller struct {
	*actor.Controller[any]
	es    *es.Store
	coord *coordinator.Controller
}

func New(esStore *es.Store) *Controller {
	ctrl := &Controller{es: esStore}
	ctrl.Controller = actor.NewController[any](func(ctx context.Context) actor.Processor[any] {
		return newProcessor(ctx, runnerModule.NewCollection(), esStore, ctrl)
	})
	esStore.AddObserver(newEventSourceObserver(ctrl))
	return ctrl
}

func (c *Controller) SetCoordinator(coord *coordinator.Controller) {
	c.coord = coord
}

func (c *Controller) coordinator() *coordinator.Controller {
	return c.coord
}

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
