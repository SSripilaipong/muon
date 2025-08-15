package runner

import (
	"context"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/server/coordinator"
	"github.com/SSripilaipong/muon/server/eventsource"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
)

type Controller struct {
	*actor.Controller[any]
}

func New(esCtrl *eventsource.Controller, coord *coordinator.Controller) *Controller {
	ctrl := &Controller{
		Controller: actor.NewController[any](func(ctx context.Context) actor.Processor[any] {
			return newProcessor(ctx, runnerModule.NewCollection(), coord)
		}),
	}
	esCtrl.AddObserver(newEventSourceObserver(ctrl))
	return ctrl
}
