package runner

import (
	"context"

	"github.com/SSripilaipong/muon/common/actor"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
)

type Controller = actor.Controller[any]

func New(remotes []RemoteNode) *Controller {
	return actor.NewController[any](func(ctx context.Context) actor.Processor[any] {
		logStore := newEventLog()
		coord := newCoordinator(logStore, remotes)
		return newProcessor(ctx, runnerModule.NewCollection(), coord)
	})
}
