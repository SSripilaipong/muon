package runner

import (
	"context"
	"fmt"

	"github.com/SSripilaipong/muon/common/actor"
	"github.com/SSripilaipong/muon/server/eventsource"
	runnerModule "github.com/SSripilaipong/muon/server/runner/module"
)

type Controller struct {
	eventCtrl *eventsource.Controller

	cancel context.CancelFunc
	done   <-chan struct{}
	msgBox chan any
}

func New(eventCtrl *eventsource.Controller) *Controller {
	return &Controller{eventCtrl: eventCtrl}
}

func (c *Controller) Start() error {
	if c.done != nil {
		return fmt.Errorf("runner has already been started")
	}

	msgBox := make(chan any)
	ctx, cancel := context.WithCancel(context.Background())
	done := actor.StartLoop(ctx, msgBox, newProcessor(ctx, runnerModule.NewCollection(), c.eventCtrl))

	c.msgBox = msgBox
	c.cancel = cancel
	c.done = done
	return nil
}

func (c *Controller) Stop() error {
	c.cancel()
	<-c.Done()
	return nil
}

func (c *Controller) Done() <-chan struct{} {
	return c.done
}
