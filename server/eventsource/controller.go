package eventsource

import (
	"github.com/SSripilaipong/muon/common/actor"
)

type Controller struct {
	*actor.Controller[any]
}

func NewController() *Controller {
	return &Controller{
		Controller: actor.NewController[any](newProcessor),
	}
}
