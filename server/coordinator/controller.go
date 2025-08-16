package coordinator

import (
	"context"

	"github.com/SSripilaipong/muon/common/actor"
)

type Controller struct {
	*actor.Controller[any]
	local Node
}

func New(local LocalNode) *Controller {
	return &Controller{
		Controller: actor.NewController[any](func(ctx context.Context) actor.Processor[any] {
			return newProcessor(ctx, local)
		}),
		local: local,
	}
}
