package actor

import (
	"context"
	"fmt"
)

type Controller[T any] struct {
	init func(ctx context.Context) Processor[T]

	cancel context.CancelFunc
	done   <-chan struct{}
	msgBox chan T
}

func NewController[T any](init func(ctx context.Context) Processor[T]) *Controller[T] {
	return &Controller[T]{init: init}
}

func (c *Controller[T]) Ch() chan<- T {
	return c.msgBox
}

func (c *Controller[T]) Start() error {
	if c.done != nil {
		return fmt.Errorf("runner has already been started")
	}

	msgBox := make(chan T)
	ctx, cancel := context.WithCancel(context.Background())
	done := StartLoop[T](ctx, msgBox, c.init(ctx))

	c.msgBox = msgBox
	c.cancel = cancel
	c.done = done
	return nil
}

func (c *Controller[T]) Stop() error {
	c.cancel()
	<-c.Done()
	return nil
}

func (c *Controller[T]) Done() <-chan struct{} {
	return c.done
}
