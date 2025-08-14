package runner

import (
	"context"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/server/eventsource"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
)

type Controller struct {
	*actor.Controller[any]
}

func New(eventCtrl *eventsource.Controller) *Controller {
	ctrl := &Controller{
		Controller: actor.NewController[any](func(ctx context.Context) actor.Processor[any] {
			return newProcessor(ctx, runnerModule.NewCollection(), eventCtrl)
		}),
	}
	eventCtrl.AddObserver(newEventSourceObserver(ctrl))
	return ctrl
}
