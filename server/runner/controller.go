package runner

import (
	"context"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/server/coordinator"
	es "github.com/SSripilaipong/muon/server/eventsource"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
)

type Controller struct {
	*actor.Controller[any]
	coord *coordinator.Controller
}

func New(esStore *es.Store) *Controller {
	ctrl := &Controller{}
	ctrl.Controller = actor.NewController[any](func(ctx context.Context) actor.Processor[any] {
		coord := ctrl.coord
		if coord == nil {
			panic("runner coordinator is not set")
		}
		return newProcessor(ctx, runnerModule.NewCollection(), esStore, coord)
	})
	esStore.AddObserver(newEventSourceObserver(ctrl))
	return ctrl
}

func (c *Controller) SetCoordinator(coord *coordinator.Controller) {
	c.coord = coord
}
