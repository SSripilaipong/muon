package eventsource

import (
	"context"

	"github.com/SSripilaipong/muon/common/actor"
)

type Controller struct {
	*actor.Controller[any]
	observer *observeSubject
}

func NewController() *Controller {
	observer := newObserveSubject()
	return &Controller{
		Controller: actor.NewController[any](func(ctx context.Context) actor.Processor[any] {
			return newProcessor(ctx, observer)
		}),
		observer: observer,
	}
}

func (c *Controller) AddObserver(observer Observer) {
	c.observer.Attach(observer)
}
