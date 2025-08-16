package prl

import (
	"context"
	"sync"
)

func Collect[T any](ctx context.Context, fs ...func() T) <-chan T {
	result := make(chan T, len(fs))
	wg := &sync.WaitGroup{}
	for _, f := range fs {
		wg.Go(func() {
			select {
			case <-ctx.Done():
			case result <- f():
			}
		})
	}
	go func() {
		defer close(result)
		wg.Wait()
	}()
	return result
}
