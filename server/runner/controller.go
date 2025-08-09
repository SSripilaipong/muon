package runner

import (
	"context"
	"fmt"

	"github.com/SSripilaipong/muto/syntaxtree/result"

	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/server/runner/module"
)

type Controller struct {
	cancel context.CancelFunc
	done   <-chan struct{}
	msgBox chan any
}

func (c *Controller) Start() error {
	if c.done != nil {
		return fmt.Errorf("runner has already been started")
	}

	msgBox := make(chan any)
	ctx, cancel := context.WithCancel(context.Background())
	done := startRunner(ctx, msgBox)

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

func (c *Controller) Run(node result.SimplifiedNode) error {
	reply := make(chan error, 1)
	defer close(reply)

	if err := chn.SendWithTimeout[any](c.msgBox, runMessage{
		moduleVersion: module.VersionDefault,
		node:          node,
		reply:         reply,
	}, channelTimeout); err != nil {
		return fmt.Errorf("cannot connect to runner: %w", err)
	}
	return chn.ReceiveWithTimeout(reply, channelTimeout).Error()
}
