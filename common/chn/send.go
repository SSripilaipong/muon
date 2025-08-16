package chn

import (
	"context"
	"fmt"
	"time"
)

func SendNoWait[T any](c chan<- T, x T) bool {
	select {
	case c <- x:
		return true
	default:
		return false
	}
}

func SendWithTimeout[T any](c chan<- T, x T, timeout time.Duration) error {
	select {
	case c <- x:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timed out")
	}
}

func SendWithContext[T any](ctx context.Context, c chan<- T, x T) error {
	select {
	case c <- x:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
